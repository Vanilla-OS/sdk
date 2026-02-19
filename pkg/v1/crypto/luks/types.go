package luks

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

// Params is a minimal LUKS parameter set, intentionally generic.
//
// Not all fields are used by all backends.
type Params struct {
	Cipher        string
	Mode          string
	Hash          string
	VolumeKeySize uint64
	DataAlignment uint64
	DataDevice    *string
}
