#ifndef NATIVE_C
#define NATIVE_C

#include <stdio.h>
#include <dlfcn.h>
#include <unistd.h>
#include <gtk/gtk.h>
#include <string.h>
#include <libappindicator/app-indicator.h>

// typedefs for dlopen
typedef void* (*app_indicator_new_fun)(const gchar*, const gchar*, AppIndicatorCategory);
typedef void* (*app_indicator_set_status_fun)(AppIndicator*, AppIndicatorStatus);
typedef void* (*app_indicator_set_icon_fun) (AppIndicator*, const gchar*);
typedef void* (*app_indicator_set_menu_fun)(AppIndicator*, GtkMenu*);

// settings
static const char* tray_name = NULL;

// reusable things
static GtkWidget *menu = NULL;
static void *appindicator_handle = NULL;

// actual tray icons
static AppIndicator *indicator = NULL;
static GtkStatusIcon *tray_icon = NULL;

// dlopen functions
app_indicator_new_fun dl_app_indicator_new;
app_indicator_set_status_fun dl_app_indicator_set_status;
app_indicator_set_icon_fun dl_app_indicator_set_icon;
app_indicator_set_menu_fun dl_app_indicator_set_menu;

void reset_menu()
{
  if (menu != NULL) {
     gtk_widget_destroy(menu);
  }
  menu = gtk_menu_new();
  if (appindicator_handle != NULL && indicator != NULL) {
    dl_app_indicator_set_menu(indicator, GTK_MENU(menu));
  }
}

// internal wrapper for go callback
void _tray_callback(GtkMenuItem *item, gpointer user_data)
{
  tray_callback(GPOINTER_TO_INT(user_data));
}

void set_menu_item(int id, const char* title, int disabled)
{
  GList *list_item = NULL;
  for (list_item = gtk_container_get_children(GTK_CONTAINER(menu)); list_item != NULL; list_item = list_item->next) {
    if (id == GPOINTER_TO_INT(g_object_get_data(G_OBJECT(list_item->data), "item-id"))) {
      gtk_menu_item_set_label(GTK_MENU_ITEM(list_item->data), title);
      gtk_widget_set_sensitive(GTK_WIDGET(list_item->data), !disabled);
      gtk_widget_show(GTK_WIDGET(list_item->data));
      return;
    }
  }

  GtkWidget *item = NULL;
  if (title != NULL && title[0] == '\0') {
    item = gtk_separator_menu_item_new();
  } else {
    item = gtk_menu_item_new_with_label(title);
  }

  if (disabled == TRUE) {
     gtk_widget_set_sensitive(GTK_WIDGET(item), FALSE);
  }

  g_object_set_data(G_OBJECT(item), "item-id", GINT_TO_POINTER(id));
  g_signal_connect(G_OBJECT(item), "activate", G_CALLBACK(_tray_callback), GINT_TO_POINTER(id));
  gtk_menu_shell_append(GTK_MENU_SHELL(menu), item);
  gtk_widget_show(GTK_WIDGET(item));
}

void create_indicator()
{
  indicator = dl_app_indicator_new (tray_name, "Test name", APP_INDICATOR_CATEGORY_APPLICATION_STATUS);
  dl_app_indicator_set_menu (indicator, GTK_MENU(menu));
  dl_app_indicator_set_status (indicator, APP_INDICATOR_STATUS_ACTIVE);
}

void status_icon_activate(GtkStatusIcon *status_icon, guint button, guint activate_time, gpointer user_data)
{
  gtk_menu_popup(GTK_MENU(menu), NULL, NULL, NULL, NULL, 0, gtk_get_current_event_time());
}

void create_status_icon()
{
  tray_icon = gtk_status_icon_new();
  g_signal_connect(G_OBJECT(tray_icon), "popup-menu", G_CALLBACK(status_icon_activate), NULL);
  g_signal_connect(G_OBJECT(tray_icon), "activate", G_CALLBACK(_tray_callback), GINT_TO_POINTER(-1));
  gtk_status_icon_set_tooltip_text(tray_icon, tray_name);
  gtk_status_icon_set_visible(tray_icon, TRUE);
}

void set_icon(const char *iconName)
{
  go_log("Set icon");
  if (tray_icon != NULL) {
    gtk_status_icon_set_from_file(tray_icon, iconName);
  }
  if (indicator != NULL && appindicator_handle != NULL) {
    dl_app_indicator_set_icon(indicator, iconName);
  }
}

void init(const char* name, int desktop)
{
    int argc = 0;
    char *argv[] = { "" };
    gtk_init(&argc, (char***)&argv);

    tray_name = strdup(name);
    reset_menu();

    if (desktop == 2) { //GNOME
      init_gtk();
    }

    if (desktop == 3) { //Unity
      init_indicator();
    }

    if (desktop == 4) { //generic
      if (init_indicator() != 0) {
        init_gtk();
      }
    }
}

int init_gtk() {
  go_log("Using GTK\n");
  create_status_icon();
  return 0;
}

int init_indicator() {
  go_log("Using libappindicator\n");
  // check if system has libappindicator1 package
  appindicator_handle = dlopen("libappindicator.so.1", RTLD_LAZY);
  if (appindicator_handle == NULL) {
    // try libappindicator3
    appindicator_handle = dlopen("libappindicator3.so.1", RTLD_LAZY);
  }

  if (appindicator_handle != NULL) {
    dl_app_indicator_new = dlsym(appindicator_handle, "app_indicator_new");
    dl_app_indicator_set_status = dlsym(appindicator_handle, "app_indicator_set_status");
    dl_app_indicator_set_menu = dlsym(appindicator_handle, "app_indicator_set_menu");
    dl_app_indicator_set_icon = dlsym(appindicator_handle, "app_indicator_set_icon");
    create_indicator();
    return 0;
  } else {
    go_log("Failed to load libappindicator shared library (via dlopen)");
    return 1;
  }
}

void native_loop()
{
  gtk_main();
}

void exit_loop()
{
  gtk_main_quit();
  // if (appindicator_handle != NULL) {
  //   dlclose(appindicator_handle);
  // }
}


#endif // NATIVE_C
