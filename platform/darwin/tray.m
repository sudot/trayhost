#import <Cocoa/Cocoa.h>

NSMenu* appMenu;
NSStatusItem* statusItem;
NSAutoreleasePool *pool;

@interface ManageHandler : NSObject
+ (IBAction)manage:(id)sender;
@end

@implementation ManageHandler
+ (IBAction)manage:(id)sender {
    tray_callback([[sender representedObject] intValue]);
}
@end

@interface UpdateHandler : NSObject
+ (IBAction)update:(id)sender;
@end

@implementation UpdateHandler
+ (IBAction)update:(id)sender {

    int itemId = [[[sender objectAtIndex: 0] autorelease] intValue];
    NSString* manageTitle = [[sender objectAtIndex: 1] autorelease];
    BOOL enabled = ![[sender objectAtIndex: 2] boolValue];

    NSMenuItem* menuItem = [appMenu itemWithTag: itemId];

    if (menuItem == nil) {
        if ([manageTitle length] == 0) {
            menuItem = [NSMenuItem separatorItem];
        } else {
            menuItem = [[[NSMenuItem alloc] initWithTitle: manageTitle
                                action:@selector(manage:) keyEquivalent:@""] autorelease];
            [menuItem setRepresentedObject:[NSNumber numberWithInt:itemId]];
            [menuItem setTarget:[ManageHandler class]];
            [menuItem setTag: itemId];
            [menuItem setEnabled: enabled];
        }
        [appMenu addItem: menuItem];
    } else {
        [menuItem setTitle: manageTitle];
        [menuItem setEnabled: enabled];
    }

}
@end

void set_menu_item(int itemId, const char *title, int disabled) {
     
    NSArray *data = [NSArray arrayWithObjects: [[NSNumber numberWithInt: itemId] autorelease], [[NSString stringWithUTF8String: title] autorelease], [[NSNumber numberWithBool:(BOOL)disabled] autorelease], nil];
    [[NSRunLoop mainRunLoop] performSelector:@selector(update:) target:[UpdateHandler class] argument:data order:1 modes:[NSArray arrayWithObjects: NSRunLoopCommonModes, NSEventTrackingRunLoopMode, nil]];

}

void set_icon(const char *iconPth) {

    NSString *iconPath = [[NSString stringWithUTF8String: iconPth] autorelease];
    NSString *iconName = [iconPath lastPathComponent];

    NSImage *icon = [NSImage imageNamed: iconName];
    if (icon == nil) {
        icon = [[[NSImage alloc] initWithContentsOfFile: iconPath] autorelease];
        NSSize size;
        size.width = 18;
        size.height = 18;
        [icon setScalesWhenResized: YES];
        [icon setSize: size];
        [icon setName: iconName];
    }

    [statusItem setImage: icon];
}

void native_loop() {
    [NSApp run];
}

void exit_loop() {
    [NSApp stop:nil];
}

int init(const char *title) {

    pool = [NSAutoreleasePool new];

    [NSApplication sharedApplication];
    [NSApp setActivationPolicy:NSApplicationActivationPolicyProhibited];

    appMenu = [[NSMenu new] autorelease];
    [appMenu setAutoenablesItems: NO];

    statusItem = [[[NSStatusBar systemStatusBar] statusItemWithLength:NSVariableStatusItemLength] retain];
    [statusItem setMenu:appMenu];
    [statusItem setHighlightMode:NO];

    return 0;
}