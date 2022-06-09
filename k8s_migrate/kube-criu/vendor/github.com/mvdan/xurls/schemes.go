/*
Copyright (c) 2014-2020 CGCL Labs
Container_Migrate is licensed under Mulan PSL v2.
You can use this software according to the terms and conditions of the Mulan PSL v2.
You may obtain a copy of Mulan PSL v2 at:
         http://license.coscl.org.cn/MulanPSL2
THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
See the Mulan PSL v2 for more details.
*/
package xurls

// SchemesNoAuthority is a sorted list of some well-known url schemes that are
// followed by ":" instead of "://". Since these are more prone to false
// positives, we limit their matching.
var SchemesNoAuthority = []string{
	`bitcoin`, // Bitcoin
	`file`,    // Files
	`magnet`,  // Torrent magnets
	`mailto`,  // Mail
	`sms`,     // SMS
	`tel`,     // Telephone
	`xmpp`,    // XMPP
}
