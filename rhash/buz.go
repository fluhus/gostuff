// Implementation of 64-bit buzhash.

package rhash

import (
	"math/rand"
)

type BuzSeed [256]uint64

// Buz implements a buzhash rolling-hash.
// Implements [hash.Hash64].
type Buz struct {
	h    uint64 // Current state
	i    int    // Number of bytes written
	b    []byte // Buffer for subtracting old bytes
	seed *BuzSeed
}

// NewBuz returns a new rolling hash with a window size of n.
func NewBuz(n int) *Buz {
	return &Buz{b: make([]byte, n), seed: defaultBuzSeed}
}

// NewBuzWithSeed returns a new rolling hash with a window size of n,
// with the given seed. Seeds can be generated with [BuzRandomSeed].
func NewBuzWithSeed(n int, seed *BuzSeed) *Buz {
	seedCopy := &BuzSeed{}
	*seedCopy = *seed
	return &Buz{b: make([]byte, n), seed: seedCopy}
}

// BuzRandomSeed returns a random seed fot the buz64 hash.
func BuzRandomSeed() *BuzSeed {
	seed := &BuzSeed{}
	for i := range seed {
		seed[i] = rand.Uint64()
	}
	return seed
}

// WriteByte updates the hash with the given byte. Always returns nil.
func (h *Buz) WriteByte(b byte) error {
	n := len(h.b)
	i := h.i % n
	if h.i >= n { // Need to subtract an old character.
		h.h ^= shift64(h.seed[h.b[i]], (n-1)%64)
	}
	h.h = shift64(h.h, 1)
	h.h ^= h.seed[b]
	h.b[i] = b
	h.i++

	return nil
}

// Write updates the hash with the given bytes.
// Always returns len(data), nil.
func (h *Buz) Write(data []byte) (int, error) {
	for _, b := range data {
		h.WriteByte(b)
	}
	return len(data), nil
}

// Sum64 returns the current hash.
func (h *Buz) Sum64() uint64 {
	return h.h
}

// Sum32 returns the current hash.
func (h *Buz) Sum32() uint32 {
	return uint32(h.h)
}

// BlockSize returns the hash's block size, which is one.
func (h *Buz) BlockSize() int {
	return 1
}

// Reset resets the hash to its initial state.
func (h *Buz) Reset() {
	h.h = 0
	h.i = 0
}

// Size returns the number of bytes Sum will return, which is eight.
func (h *Buz) Size() int {
	return 8
}

// Sum appends the current hash to b and returns the resulting slice.
func (h *Buz) Sum(b []byte) []byte {
	s := h.Sum64()
	for range h.Size() {
		b = append(b, byte(s))
		s >>= 8
	}
	return b
}

// Returns a cyclic binary left shift by n bits.
func shift64(x uint64, n int) uint64 {
	return x<<n | x>>(64-n)
}

// Random single-byte hashes.
var defaultBuzSeed = &BuzSeed{
	0x75a494d1541fbc6, 0xd526e3fdd0bea3c7, 0xd04e3eac66b233f, 0xf6f31d3a1a7dc222,
	0xfc1a9f426c0f84e4, 0xf2ac01cb50518375, 0x500f4db1da25b019, 0x24062ea231ca55bf,
	0x52afae6d4dc824cc, 0xcc838d5a5f8970b6, 0xabcc40267f2a3806, 0xaf2939d0a84e2828,
	0x918bf35c02097e37, 0x547dac6ec5152648, 0x41f3a34be76c760c, 0xe03a50e492bc37b3,
	0x4b2dcdba9627c926, 0xe351d62282c74d66, 0x92e8579242cfe718, 0x5bef0d63dc35595e,
	0xbb6a6f9956e35194, 0x4cb95cfd0881c69d, 0x4235f12344ff932b, 0xc55d679319513a8f,
	0x7f7f6148f4dfab1d, 0x4d5dbaccc4a7b030, 0xeda326161dad7579, 0xba0fac64641178f7,
	0x3687facc0faa8614, 0x7d886bac2b333e7a, 0x6c881755d7cce7e1, 0xe35611a8712ec60e,
	0xdc17afe56925789b, 0xfecf7587fc689832, 0x94bfbbfe9082be61, 0x39100b86c5962ee3,
	0x7f07f1aec5f9c27c, 0x6f284852f62d7a92, 0xebdbff207d7452b0, 0x6f01e4f1df825e79,
	0x6ad8cb66e26873d8, 0xa60cf9e64d49b36a, 0x914fe6afce1479a7, 0x720ba589fdce07bf,
	0x1f82029d83b228f2, 0x1e1d85c50df8bff2, 0xbf095470ce998aa4, 0xb3fa4be3db7c2e95,
	0xc66c984e51cc6efa, 0xdae702fb44646c68, 0xf08f3aa4edb724ac, 0x52e2468427a1de62,
	0x24b56118b69d4701, 0x3753b13b0cd62cac, 0xe373655df4cb3ae6, 0xc48056cb98a0950c,
	0xc9c2de155bf13d3c, 0x2557b6b645c024fd, 0x642f35503d19cc25, 0xa816f75953ff8ed2,
	0xa6a3625532490cb, 0x6d8dd745853001de, 0xfcda0a887efd5146, 0x3395b472c17e93bf,
	0x7ae0693314c8422, 0xa940168417caf59c, 0x4074aa48b1246fee, 0x7cd26ef7649e1b,
	0x58847d34728e3d71, 0x3aa5159380b49b, 0x537a627406c94f23, 0x73bd1314bf06fefd,
	0x34204c0495b6c70c, 0xf3e89711bfd6b8f1, 0x700c9e56257791a, 0xeba59ebf8cee7e22,
	0xb19ac0ae4fd0e93d, 0xa61506e71c6cb458, 0xd4109a9d01b26219, 0x8273c279d96358f1,
	0x7d669589c566904c, 0x65a0eb44798347b9, 0x766cf27ab2ad6498, 0x987b9a8d452c51cc,
	0x1e1f8e33a63383cb, 0x80a06577e1f7ec76, 0x82c94646f19ab354, 0x231784142dd2ecee,
	0x134361ebc7296c94, 0x729071d106183edd, 0xe88a585faa0fde53, 0x889b9e02437a2cab,
	0xa4cfd4f43b576b7, 0xde4b12d99bf73257, 0x68357727d42ecc32, 0x7f8f622af2444da9,
	0x73509678a92b9bc8, 0x3fd39db1de54f32c, 0xf1517c4bb5ea3b99, 0xc5f1fa8aa9e89faf,
	0x4bcb3fc8f8e67efb, 0x4e88b9d1f31b5bdb, 0x45bd940bd029ec7d, 0x9992a46d521889a2,
	0xdf956604105c20b9, 0xeacb9ec11109d6fb, 0xd82889e054171908, 0x5f0c9bde49523051,
	0x11cc350cb39ae65c, 0x97e09c54b5c75c8b, 0x737112ade48157e4, 0x86a6795bf35790f2,
	0x8579db21057fee26, 0xda9d97b930ea67f7, 0x6c0176c93d2f27ac, 0xef082a73a8316e07,
	0x8db33d5c9bf515ae, 0x341c24b67b64152a, 0xba0a0b7ef2ad506d, 0xb3b1eed3fc4dcb37,
	0x76b14cf89ed28d1, 0x230abb2c88e16a1, 0x6d256ef7941b71d0, 0x5a6b4231726b9bef,
	0x7a868af7e61b3df, 0x78a5b10a47eda84c, 0x522ee19558a1e1f7, 0x412188c4b9cd8633,
	0x7697ef7a73cbc019, 0x36fe53556d165f8, 0x32f6c47cbf609c97, 0x6c04f6a34b390aa7,
	0xc9ef5d2ebc812480, 0xfb89bf4f13d8eccd, 0x82f780ac7e3c408d, 0x884696dae93e59c0,
	0x70e13f2c0281c6f6, 0x103c8008f0331d2c, 0xfa6fbd2119dcf1e1, 0x4f4a7b5e616e4a8d,
	0xb61097b1847c26f9, 0xff657af17504d685, 0x62e6ea13cd2d39f7, 0x99f5e05619e209d4,
	0xa6079200acba884d, 0x18a8de9e8eac3758, 0x79ead3140eeb28a4, 0x9eb43eb28d8f588a,
	0x8d693728923aa9f5, 0xf3d3fedea7c3b1c5, 0xcf8113e98da03a39, 0x9cea4c5eaf28d276,
	0x27ff4ea1f5f4e86a, 0x8efeafad91f45573, 0xa4ce200d57829b0a, 0x21eeec5210d6cb74,
	0xdc5198a841ece72b, 0x983b90fe1d650b06, 0xca8324ab7ded3e0d, 0xf3b30090fb615fe5,
	0x62972af74cff4da2, 0xa5c7ccb0640115c9, 0x9486624d1c15c3ec, 0x1c9b2eafdb861d77,
	0x2bd32a36d3e53740, 0x678f5184fd5665d, 0x497ae37164bb7af, 0x786657209be29aa1,
	0xcf2df88f9d4bd2c6, 0x345c188d10489e2, 0x971f81ee6f7cbaeb, 0x4072ff5164578516,
	0xad888677630a20b5, 0x6fd79c924b66500e, 0x1a16d63320a639b7, 0xc09b62dc19924b93,
	0x6fc05884adf20a5c, 0x6dda12167e884822, 0xd45591010b5ee8d5, 0x49583521b1af62cb,
	0x61ccad131d5d4e93, 0xdfa7bdae833942f0, 0x1e111658435f1b0e, 0xe214dd825d3a3f68,
	0xdb48374a0b61c0c1, 0x72e6f8439bc73df2, 0xe09600a0007426c1, 0x9d9161193dffb480,
	0x10b24a104bcc0387, 0x1c529c3f8209fef4, 0x68b5ae930cecc8ea, 0x1fbab1e973e16f63,
	0x3dc73f2e478a16c6, 0x9c9515d5975b62fb, 0xf36827d8cc0c7fed, 0x336511a0a7d524ca,
	0x94b397ab51409bc2, 0xdb1d1f6da07a589a, 0xbe504b7b2968fb89, 0x877f2d757115ac1a,
	0x5e3b99683df2efe4, 0x26413850ed6fa805, 0xd233def75ec0321, 0xbe207fd9dc5c3a87,
	0xe697da338e89e302, 0x4af84ce0c8b4f7c8, 0x695cba2f5faa46f1, 0x67a2aa311565960a,
	0xed99a9c07b51bd1f, 0x22c81c1ea975da80, 0xcdb1ce1ed9705dbe, 0xc9f6312b829e4648,
	0xf00b1b7614cd7881, 0xfe9739e7f46805ea, 0x31c564faa42c3ea5, 0xeba60419e23fca36,
	0xc3876743d671a7ae, 0x18fe2808cb474033, 0x5de26540cad04df7, 0x97d9f509567e7c75,
	0x22fc9c0b56a2ca67, 0x58cdab8e0948c2e, 0x8449025deca2f31e, 0x1f86b4670a04a485,
	0xb2635acb2f9c8400, 0xc304fd87df987a77, 0xb38184573350f1b, 0x96498cd171469702,
	0x75fbfea8ff2cbe54, 0x57d6fbdb7677ce2e, 0x7f728e585b9f36e8, 0x929593b474a6be51,
	0x81b829ee481749a8, 0xf73fd958496584e, 0x4d62aec3c157d2f8, 0xf0117fd0e37966f6,
	0xa72aecb159507b5, 0x31133ed09659c5d4, 0x63ea919e0afa37a, 0x287423219bc878ed,
	0x70625dc642746b20, 0x80dff18878ea68e7, 0x222ae35ddc4cc16c, 0x672cab74d0fbfd82,
	0xa22f016f8324275a, 0x7c04393534369e41, 0xb36bfdb389f6aab, 0xa2c99d0581f0b1d5,
	0x251530db20f0ff7e, 0xe78c7cbf708fc01d, 0xaf0528f43781f369, 0x8e6899ddaaa18643,
	0x4ce0021f2bb60c2, 0xbed555bf9d82a5, 0xf94748a963109133, 0xb5a13fc33de246ce,
}
