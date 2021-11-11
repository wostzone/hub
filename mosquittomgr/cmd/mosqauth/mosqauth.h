/* Code generated by cmd/cgo; DO NOT EDIT. */

/* package github.com/wostzone/hub/mosquittomgr/cmd/mosqauth */


#line 1 "cgo-builtin-export-prolog"

#include <stddef.h> /* for ptrdiff_t below */

#ifndef GO_CGO_EXPORT_PROLOGUE_H
#define GO_CGO_EXPORT_PROLOGUE_H

#ifndef GO_CGO_GOSTRING_TYPEDEF
typedef struct { const char *p; ptrdiff_t n; } _GoString_;
#endif

#endif

/* Start of preamble from import "C" comments.  */




/* End of preamble from import "C" comments.  */


/* Start of boilerplate cgo prologue.  */
#line 1 "cgo-gcc-export-header-prolog"

#ifndef GO_CGO_PROLOGUE_H
#define GO_CGO_PROLOGUE_H

typedef signed char GoInt8;
typedef unsigned char GoUint8;
typedef short GoInt16;
typedef unsigned short GoUint16;
typedef int GoInt32;
typedef unsigned int GoUint32;
typedef long long GoInt64;
typedef unsigned long long GoUint64;
typedef GoInt64 GoInt;
typedef GoUint64 GoUint;
typedef __SIZE_TYPE__ GoUintptr;
typedef float GoFloat32;
typedef double GoFloat64;
typedef float _Complex GoComplex64;
typedef double _Complex GoComplex128;

/*
  static assertion to make sure the file is being used on architecture
  at least with matching size of GoInt.
*/
typedef char _check_for_64_bit_pointer_matching_GoInt[sizeof(void*)==64/8 ? 1:-1];

#ifndef GO_CGO_GOSTRING_TYPEDEF
typedef _GoString_ GoString;
#endif
typedef void *GoMap;
typedef void *GoChan;
typedef struct { void *t; void *v; } GoInterface;
typedef struct { void *data; GoInt len; GoInt cap; } GoSlice;

#endif

/* End of boilerplate cgo prologue.  */

#ifdef __cplusplus
extern "C" {
#endif

extern void AuthPluginInit(GoSlice keys, GoSlice values, GoInt authOptsNum);

// AuthUnpwdCheck checks for a correct username/password
// This matches the given password against the stored password hash
//  clientID used to connect
//  username is the login user name
//  password is the login password
//  clientIP
//  certSubjName when authenticated using a certificate instead of username/password
// Returns:
//  MOSQ_ERR_SUCCESS if the user is authenticated
//  MOSQ_ERR_PLUGIN_DEFER if we do not wish to handle this check
extern GoUint8 AuthUnpwdCheck(GoString clientID, GoString username, GoString password, GoString clientIP, GoString certSubjName);

// AuthAclCheck checks if the user has access to the topic
// If certificate authentication was used the certSubjName includes the OU of the client.
// The authorizer engine can decide to give extra access to clients based on their OU.
//
// This:
//   1. determines the thingID to access from the topic
//   2. determine the groups the Thing is in
//   3. determine the highest permission of the user if a member of one of those groups
//
//  clientID used to connect to the message bus
//  userID login ID of user when logging in with username/password
//  certSubjName: certificate subject when client certificate authentication is used
//  topic to validate
//  access: MOSQ_ACL_SUBSCRIBE, MOSQ_ACL_READ, MOSQ_ACL_WRITE
//
// returns
//  MOSQ_ERR_ACL_DENIED if access was not granted
//  MOSQ_ERR_UNKNOWN for an application specific error
//  MOSQ_ERR_SUCCESS if access is granted
//  MOSQ_ERR_PLUGIN_DEFER if we do not wish to handle this check
extern GoUint8 AuthAclCheck(GoString clientID, GoString userID, GoString certSubjName, GoString topic, GoInt access);
extern void AuthPluginCleanup();

#ifdef __cplusplus
}
#endif
