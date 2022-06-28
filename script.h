#import <Foundation/Foundation.h>

NSAppleScript* compileScript(char* script);
NSAppleEventDescriptor* executeScript(NSAppleScript* script);
const char* getStringFromDescriptor(NSAppleEventDescriptor* descriptor, int index);
int getIntFromDescriptor(NSAppleEventDescriptor* descriptor, int index);
float getFloatFromDescriptor(NSAppleEventDescriptor* descriptor, int index);