#include <stdio.h>
#include <windows.h>
#include <shellapi.h>
#include <strsafe.h>
#include <tchar.h>

#define WM_MYMESSAGE (WM_USER + 1)
#define TRAYHOST_ICON_ID 100
#define MAX_LOADSTRING 100

// HINSTANCE hInst;
HWND hWnd;
HMENU hSubMenu;
TCHAR szTitle[MAX_LOADSTRING];
TCHAR szWindowClass[MAX_LOADSTRING];
// wchar_t *title;
NOTIFYICONDATA nid;

ATOM                MyRegisterClass(HINSTANCE hInstance);
HWND                InitInstance(HINSTANCE, int);
LRESULT CALLBACK    WndProc(HWND, UINT, WPARAM, LPARAM);

extern void go_log(char*);
extern void tray_callback(int itemId);

void reset_menu()
{
    if (hSubMenu != NULL) {
        DestroyMenu(hSubMenu);
    }
    hSubMenu = CreatePopupMenu();
}

void set_menu_item(int id, const char* go_title, int disabled)
{
    LPTSTR title = (LPTSTR)go_title;

    MENUITEMINFOW menu_item_info;
    memset(&menu_item_info, 0, sizeof(MENUITEMINFO));
    menu_item_info.cbSize = sizeof(MENUITEMINFO);
    
    BOOL alreadyExists = GetMenuItemInfoW(hSubMenu, id, FALSE, &menu_item_info);
    menu_item_info.fMask = MIIM_STRING | MIIM_DATA | MIIM_FTYPE;

    if (title != NULL && title[0] == '\0') {
        menu_item_info.fType = MFT_SEPARATOR;
    } else {
        menu_item_info.fType = MFT_STRING;
    }

    if (disabled == TRUE) {
        menu_item_info.fMask = menu_item_info.fMask | MIIM_STATE;
        menu_item_info.fState = MFS_GRAYED;
    }
    
    menu_item_info.dwTypeData = title;

    if (alreadyExists == TRUE) {
        SetMenuItemInfo(hSubMenu, id, FALSE, &menu_item_info);
    } else {
        menu_item_info.fMask = menu_item_info.fMask | MIIM_ID;
        menu_item_info.wID = id;
        InsertMenuItem(hSubMenu, id, FALSE, &menu_item_info);
    }
}

void set_icon(const char *go_iconPth)
{
    LPCTSTR iconPth = (LPCTSTR)go_iconPth;

    NOTIFYICONDATA nid;
    memset(&nid, 0, sizeof(NOTIFYICONDATA));
    nid.cbSize = sizeof(NOTIFYICONDATA);
    nid.uID = TRAYHOST_ICON_ID;
    nid.hWnd = hWnd;
    nid.hIcon = LoadImage(NULL, iconPth, IMAGE_ICON, 64, 64, LR_LOADFROMFILE);
    nid.uFlags = NIF_ICON;

    Shell_NotifyIcon(NIM_MODIFY, &nid);
}

void native_loop()
{
    MSG msg;
    // Main message loop:
    while (GetMessage(&msg, NULL, 0, 0))
    {
        TranslateMessage(&msg);
        DispatchMessage(&msg);
    }
}

void init(const char *go_title)
{
    LPCTSTR title = (LPCTSTR)go_title;
    HINSTANCE hInstance = GetModuleHandle(NULL);

    StringCchCopy(szTitle, MAX_LOADSTRING, title);
    StringCchCopy(szWindowClass, MAX_LOADSTRING, L"MyClass");
    MyRegisterClass(hInstance);

    hWnd = InitInstance(hInstance, FALSE); // Don't show window
    if (!hWnd)
    {
        return;
    }

    NOTIFYICONDATA nid;
    memset(&nid, 0, sizeof(NOTIFYICONDATA));
    nid.cbSize = sizeof(NOTIFYICONDATA);
    StringCchCopy(nid.szTip, 64, szTitle);
    nid.uID = TRAYHOST_ICON_ID;
    nid.hWnd = hWnd;
    nid.uCallbackMessage = WM_MYMESSAGE; 
    nid.uFlags = NIF_TIP | NIF_MESSAGE;

    go_log("add_icon");
    Shell_NotifyIcon(NIM_ADD, &nid);
    reset_menu();
}

void exit_loop() {
    Shell_NotifyIcon(NIM_DELETE, &nid);
    PostQuitMessage(0);
}


ATOM MyRegisterClass(HINSTANCE hInstance)
{
    WNDCLASSEX wcex;

    wcex.cbSize = sizeof(WNDCLASSEX);

    wcex.style          = CS_HREDRAW | CS_VREDRAW;
    wcex.lpfnWndProc    = WndProc;
    wcex.cbClsExtra     = 0;
    wcex.cbWndExtra     = 0;
    wcex.hInstance      = hInstance;
    wcex.hIcon          = LoadIcon(NULL, IDI_APPLICATION);
    wcex.hCursor        = LoadCursor(NULL, IDC_ARROW);
    wcex.hbrBackground  = (HBRUSH)(COLOR_WINDOW+1);
    wcex.lpszMenuName   = 0;
    wcex.lpszClassName  = szWindowClass;
    wcex.hIconSm        = LoadIcon(NULL, IDI_APPLICATION);

    return RegisterClassEx(&wcex);
}

HWND InitInstance(HINSTANCE hInstance, int nCmdShow)
{

    hWnd = CreateWindow(szWindowClass, szTitle, WS_OVERLAPPEDWINDOW,
                        CW_USEDEFAULT, 0, CW_USEDEFAULT, 0, NULL, NULL, hInstance, NULL);

    if (!hWnd)
    {
        return 0;
    }

    ShowWindow(hWnd, nCmdShow);
    UpdateWindow(hWnd);

    return hWnd;
}

void ShowMenu(HWND hWnd)
{
    POINT p;
    GetCursorPos(&p);
    SetForegroundWindow(hWnd); // Win32 bug work-around
    TrackPopupMenu(hSubMenu, TPM_BOTTOMALIGN | TPM_LEFTALIGN, p.x, p.y, 0, hWnd, NULL);

}

LRESULT CALLBACK WndProc(HWND hWnd, UINT message, WPARAM wParam, LPARAM lParam)
{
    switch (message)
    {
        case WM_COMMAND:
            tray_callback(LOWORD(wParam));
            break;
        case WM_DESTROY:
            exit_loop();
            break;
        case WM_MYMESSAGE:
            switch(lParam)
            {
                case WM_RBUTTONUP:
                    ShowMenu(hWnd);
                    break;
                case WM_LBUTTONUP:
                    tray_callback(-1);
                    break;
                default:
                    return DefWindowProc(hWnd, message, wParam, lParam);
            };
            break;
        default:
            return DefWindowProc(hWnd, message, wParam, lParam);
    }
    return 0;
}
