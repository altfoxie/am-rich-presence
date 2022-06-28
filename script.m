#import "script.h"

// fucking shitcode idk apple sdks and objc

NSAppleScript* compileScript(char* script) {
    NSString* scriptString = [NSString stringWithUTF8String:script];
    NSAppleScript* appleScript = [[NSAppleScript alloc] initWithSource:scriptString];
    if (![appleScript compileAndReturnError:nil]) return nil;
    return appleScript;
}

NSAppleEventDescriptor* executeScript(NSAppleScript* script) {
    return [script executeAndReturnError:nil];
}

const char* getStringFromDescriptor(NSAppleEventDescriptor* descriptor, int index) {
    return [[descriptor descriptorAtIndex:index] stringValue].UTF8String;
}

int getIntFromDescriptor(NSAppleEventDescriptor* descriptor, int index) {
    return [[descriptor descriptorAtIndex:index] int32Value];
}

float getFloatFromDescriptor(NSAppleEventDescriptor* descriptor, int index) {
    return [[descriptor descriptorAtIndex:index] floatValue];
}