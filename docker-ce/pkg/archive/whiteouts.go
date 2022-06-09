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
package archive

// Whiteouts are files with a special meaning for the layered filesystem.
// Docker uses AUFS whiteout files inside exported archives. In other
// filesystems these files are generated/handled on tar creation/extraction.

// WhiteoutPrefix prefix means file is a whiteout. If this is followed by a
// filename this means that file has been removed from the base layer.
const WhiteoutPrefix = ".wh."

// WhiteoutMetaPrefix prefix means whiteout has a special meaning and is not
// for removing an actual file. Normally these files are excluded from exported
// archives.
const WhiteoutMetaPrefix = WhiteoutPrefix + WhiteoutPrefix

// WhiteoutLinkDir is a directory AUFS uses for storing hardlink links to other
// layers. Normally these should not go into exported archives and all changed
// hardlinks should be copied to the top layer.
const WhiteoutLinkDir = WhiteoutMetaPrefix + "plnk"

// WhiteoutOpaqueDir file means directory has been made opaque - meaning
// readdir calls to this directory do not follow to lower layers.
const WhiteoutOpaqueDir = WhiteoutMetaPrefix + ".opq"
