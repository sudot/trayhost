
// #cgo CFLAGS: -DDARWIN -x objective-c -fobjc-no-arc
// #cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

NSMenu* appMenu;
NSAutoreleasePool *pool;

extern void tray_callback(int itemId);

@interface ManageHandler : NSObject
+ (IBAction)manage:(id)sender;
@end

@implementation ManageHandler
+ (IBAction)manage:(id)sender {
    tray_callback([[sender representedObject] intValue]);
}
@end

@interface UpdateHandler : NSObject
+ (IBAction)test:(id)sender;
@end

@implementation UpdateHandler
+ (IBAction)test:(id)sender {
    // NSLog(@"test %@", sender);

    int itemId = [[sender objectAtIndex: 0] intValue];
    NSString* manageTitle = [sender objectAtIndex: 1];
    BOOL enabled = (BOOL)[sender objectAtIndex: 2];

    NSMenuItem* menuItem = [appMenu itemWithTag: itemId];

    if (menuItem == nil) {
        menuItem = [[[NSMenuItem alloc] initWithTitle: manageTitle
                                action:@selector(manage:) keyEquivalent:@""] autorelease];
        NSLog(@"Create item %@", menuItem);
        [menuItem setRepresentedObject:[NSNumber numberWithInt:itemId]];
        [menuItem setTarget:[ManageHandler class]];
        [menuItem setTag: itemId];
        [menuItem setEnabled: enabled];
        [appMenu addItem: menuItem];
    } else {
        NSLog(@"Update item %@", menuItem);
        [menuItem setTitle: manageTitle];
        [menuItem setEnabled: enabled];
    }

}
@end

void add_menu_item(int itemId, const char *title, int disabled) {

    // [UpdateHandler test:nil];
     
    NSArray *data = [NSArray arrayWithObjects: [[NSNumber numberWithInt: itemId] autorelease], [[NSString stringWithUTF8String: title] autorelease], [[NSNumber numberWithBool:(BOOL)disabled] autorelease], nil];
    [[NSRunLoop mainRunLoop] performSelector:@selector(test:) target:[UpdateHandler class] argument:data order:1 modes:[NSArray arrayWithObjects: NSRunLoopCommonModes, NSEventTrackingRunLoopMode, nil]];

}

void native_loop() {
    [NSApp run];
}

void exit_loop() {
    [NSApp stop:nil];
    // [pool release];
}

int init(const char *title, unsigned char imageDataBytes[], unsigned int imageDataLen) {

    pool = [NSAutoreleasePool new];
    // pool = [[NSAutoreleasePool alloc] init];

    [NSApplication sharedApplication];
    [NSApp setActivationPolicy:NSApplicationActivationPolicyProhibited];

    appMenu = [[NSMenu new] autorelease];
    [appMenu setAutoenablesItems: NO];

    NSData *iconData = [NSData dataWithBytes:imageDataBytes length:imageDataLen];
    NSImage *icon = [[NSImage alloc] initWithData:iconData];

    NSStatusItem* statusItem = [[[NSStatusBar systemStatusBar] statusItemWithLength:NSVariableStatusItemLength] retain];
    [statusItem setMenu:appMenu];
    [statusItem setImage:icon];
    [statusItem setHighlightMode:NO];

    return 0;
}