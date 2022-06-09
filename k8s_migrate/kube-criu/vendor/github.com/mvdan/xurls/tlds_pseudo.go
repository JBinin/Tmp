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

// PseudoTLDs is a sorted list of some widely used unofficial TLDs.
//
// Sources:
//  * https://en.wikipedia.org/wiki/Pseudo-top-level_domain
//  * https://en.wikipedia.org/wiki/Category:Pseudo-top-level_domains
//  * https://tools.ietf.org/html/draft-grothoff-iesg-special-use-p2p-names-00
//  * https://www.iana.org/assignments/special-use-domain-names/special-use-domain-names.xhtml
var PseudoTLDs = []string{
	`bit`,       // Namecoin
	`example`,   // Example domain
	`exit`,      // Tor exit node
	`gnu`,       // GNS by public key
	`i2p`,       // I2P network
	`invalid`,   // Invalid domain
	`local`,     // Local network
	`localhost`, // Local network
	`onion`,     // Tor hidden services
	`test`,      // Test domain
	`zkey`,      // GNS domain name
}
