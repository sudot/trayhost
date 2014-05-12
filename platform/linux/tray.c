#ifndef NATIVE_C
#define NATIVE_C

#include <stdio.h>
#include <dlfcn.h>
#include <unistd.h>
#include <gtk/gtk.h>
#include <gio/gio.h>
#include <gdk-pixbuf/gdk-pixbuf.h>
#include <string.h>
#include <libappindicator/app-indicator.h>

char *icon = NULL;
size_t iconSize = 0;
const char *menu_title = NULL;
const char *url = NULL;
GtkWidget *menu = NULL;
void *appindicator_handle = NULL;
char tmpIconNameBuf[32];
AppIndicator *indicator = NULL;

typedef void* (*app_indicator_new_fun)(const gchar*, const gchar*, AppIndicatorCategory);
typedef void* (*app_indicator_set_status_fun)(AppIndicator*, AppIndicatorStatus);
typedef void* (*app_indicator_set_attention_icon_full_fun) (AppIndicator*,  const gchar* ,const gchar*);
typedef void* (*app_indicator_set_menu_fun)(AppIndicator*,GtkMenu*);

// implemented in go
extern void tray_callback(int itemId);

void reset_menu() {
  if (menu != NULL) {
     gtk_widget_destroy(menu);
  }
  menu = gtk_menu_new();
  if (appindicator_handle != NULL && indicator != NULL) {
    app_indicator_set_menu_fun app_indicator_set_menu;
    app_indicator_set_menu = dlsym(appindicator_handle, "app_indicator_set_menu");
    app_indicator_set_menu(indicator, GTK_MENU(menu));
  }
}

// internal wrapper for go callback
void _tray_callback(GtkMenuItem *item, gpointer user_data)
{
  tray_callback(GPOINTER_TO_INT(user_data));
  gpointer data = g_object_get_data(G_OBJECT(item), "item-id");
}

void add_menu_item(int id, const char* title, int disabled) {
  GList *list_item = NULL;
  for (list_item = gtk_container_get_children(GTK_CONTAINER(menu)); list_item != NULL; list_item = list_item->next) {
    if (id == GPOINTER_TO_INT(g_object_get_data(G_OBJECT(list_item->data), "item-id"))) {
      gtk_menu_item_set_label(GTK_MENU_ITEM(list_item->data), title);
      if (disabled == TRUE) {
        gtk_widget_set_sensitive(GTK_WIDGET(list_item->data), FALSE);
      }
      return;
    }
  }

  GtkWidget *item = NULL;
  if (title == "") {
    item = gtk_separator_menu_item_new();
  } else {
    item = gtk_menu_item_new_with_label(title);
  }
  
  if (disabled == TRUE) {
     gtk_widget_set_sensitive(GTK_WIDGET(item), FALSE);
  }
  g_object_set_data(G_OBJECT(item), "item-id", GINT_TO_POINTER(id));
  g_signal_connect(G_OBJECT(item), "activate", G_CALLBACK(_tray_callback), GINT_TO_POINTER(id));
  gtk_widget_show(item);
  gtk_menu_shell_append(GTK_MENU_SHELL(menu), item);
}

void create_indicator(void *appindicator_handle)
{
  app_indicator_new_fun app_indicator_new;
  app_indicator_set_status_fun app_indicator_set_status;
  app_indicator_set_menu_fun app_indicator_set_menu;

  app_indicator_new = dlsym(appindicator_handle, "app_indicator_new");
  app_indicator_set_status = dlsym(appindicator_handle, "app_indicator_set_status");
  app_indicator_set_menu = dlsym(appindicator_handle, "app_indicator_set_menu");

  // write icon to temp file, otherwise imposible to set in libappindicator
  int fd = -1;
  memset(tmpIconNameBuf, 0, sizeof(tmpIconNameBuf));
  strncpy(tmpIconNameBuf,"/tmp/storageguiicon-XXXXXX",26);
  fd = mkstemp(tmpIconNameBuf);

  if (fd > 0) {
    if(write(fd, icon, iconSize) == -1) {
      fprintf(stderr, "Failed to write icon data into temp file\n");
    }
  } else {
    fprintf(stderr, "Failed to create temp file for icon\n");
  }

  indicator = app_indicator_new (menu_title,
                                 tmpIconNameBuf,
                                 APP_INDICATOR_CATEGORY_APPLICATION_STATUS);

  app_indicator_set_menu (indicator, GTK_MENU(menu));
  app_indicator_set_status (indicator, APP_INDICATOR_STATUS_ACTIVE);
}

static void status_icon_activate(GtkStatusIcon *status_icon, guint button, guint activate_time, gpointer user_data)
{
    gtk_menu_popup(GTK_MENU(menu), NULL, NULL, NULL, NULL, 0, gtk_get_current_event_time());
}

void create_status_icon()
{
    GError *error = NULL;
    GInputStream *stream = g_memory_input_stream_new_from_data(icon, iconSize, NULL);
    GdkPixbuf *pixbuf = gdk_pixbuf_new_from_stream(stream, NULL, &error);
    if (error)
        fprintf(stderr, "Unable to create PixBuf: %s\n", error->message);

    GtkStatusIcon *tray_icon = gtk_status_icon_new_from_pixbuf(pixbuf);

    g_signal_connect(G_OBJECT(tray_icon), "popup-menu", G_CALLBACK(status_icon_activate), NULL);
    gtk_status_icon_set_tooltip_text(tray_icon, menu_title);
    gtk_status_icon_set_visible(tray_icon, TRUE);
}

void init(const char* title, unsigned char *imageData, unsigned int imageDataLen)
{
    int argc = 0;
    char *argv[] = { "" };
    gtk_init(&argc, (char***)&argv);

    menu_title = title;
    icon = imageData;
    iconSize = imageDataLen;
    reset_menu();

    // check if system has libappindicator1 package
    appindicator_handle = dlopen("libappindicator3.so.1", RTLD_LAZY);
    if(appindicator_handle == NULL) {
      create_status_icon();
    } else {
      create_indicator(appindicator_handle);
    }
}

void native_loop()
{
  gtk_main();
}

void exit_loop()
{
  gtk_main_quit();
}


#endif // NATIVE_C
