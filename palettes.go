package main

var c64palettes = map[string][16]colorInfo{
	"vice": {
		{0, RGB{0x00, 0x00, 0x00}},
		{1, RGB{0xff, 0xff, 0xff}},
		{2, RGB{0xbc, 0x52, 0x41}},
		{3, RGB{0x8f, 0xef, 0xfb}},
		{4, RGB{0xb9, 0x56, 0xeb}},
		{5, RGB{0x7e, 0xdb, 0x40}},
		{6, RGB{0x55, 0x3f, 0xe4}},
		{7, RGB{0xff, 0xff, 0x77}},
		{8, RGB{0xc1, 0x7b, 0x1d}},
		{9, RGB{0x82, 0x63, 0x00}},
		{10, RGB{0xf4, 0x94, 0x86}},
		{11, RGB{0x72, 0x72, 0x72}},
		{12, RGB{0xa4, 0xa4, 0xa4}},
		{13, RGB{0xcd, 0xff, 0x98}},
		{14, RGB{0x9e, 0x8d, 0xff}},
		{15, RGB{0xd5, 0xd5, 0xd5}},
	},
	"vice old lum": {
		{0, RGB{0x00, 0x00, 0x00}},
		{1, RGB{0xff, 0xff, 0xff}},
		{2, RGB{0xa9, 0x38, 0x26}},
		{3, RGB{0xae, 0xff, 0xff}},
		{4, RGB{0xdf, 0x82, 0xff}},
		{5, RGB{0x7e, 0xdb, 0x40}},
		{6, RGB{0x55, 0x3f, 0xe4}},
		{7, RGB{0xf7, 0xff, 0x6d}},
		{8, RGB{0xe7, 0xa4, 0x53}},
		{9, RGB{0x82, 0x63, 0x00}},
		{10, RGB{0xf4, 0x94, 0x86}},
		{11, RGB{0x5c, 0x5c, 0x5c}},
		{12, RGB{0xb0, 0xb0, 0xb0}},
		{13, RGB{0xc4, 0xff, 0x8f}},
		{14, RGB{0xaa, 0x99, 0xff}},
		{15, RGB{0xf2, 0xf2, 0xf2}},
	},
	"pepto": {
		{0, RGB{0x00, 0x00, 0x00}},
		{1, RGB{0xff, 0xff, 0xff}},
		{2, RGB{0x68, 0x37, 0x2b}},
		{3, RGB{0x70, 0xa4, 0xb2}},
		{4, RGB{0x6f, 0x3d, 0x86}},
		{5, RGB{0x58, 0x8d, 0x43}},
		{6, RGB{0x35, 0x28, 0x79}},
		{7, RGB{0xb8, 0xc7, 0x6f}},
		{8, RGB{0x6f, 0x4f, 0x25}},
		{9, RGB{0x43, 0x39, 0x00}},
		{10, RGB{0x9a, 0x67, 0x59}},
		{11, RGB{0x44, 0x44, 0x44}},
		{12, RGB{0x6c, 0x6c, 0x6c}},
		{13, RGB{0x9a, 0xd2, 0x84}},
		{14, RGB{0x6c, 0x5e, 0xb5}},
		{15, RGB{0x95, 0x95, 0x95}},
	},
	"pantaloon": {
		{0, RGB{0x00, 0x00, 0x00}},
		{1, RGB{0xff, 0xff, 0xff}},
		{2, RGB{0x68, 0x37, 0x2b}},
		{3, RGB{0x83, 0xf0, 0xdc}},
		{4, RGB{0x6f, 0x3d, 0x86}},
		{5, RGB{0x59, 0xcd, 0x36}},
		{6, RGB{0x41, 0x37, 0xcd}},
		{7, RGB{0xb8, 0xc7, 0x6f}},
		{8, RGB{0xd1, 0x7f, 0x30}},
		{9, RGB{0x43, 0x39, 0x00}},
		{10, RGB{0x9a, 0x67, 0x59}},
		{11, RGB{0x5b, 0x5b, 0x5b}},
		{12, RGB{0x8e, 0x8e, 0x8e}},
		{13, RGB{0x9d, 0xff, 0x9d}},
		{14, RGB{0x75, 0xa1, 0xec}},
		{15, RGB{0xc1, 0xc1, 0xc1}},
	},
	"archmage": {
		{0, RGB{0x00, 0x00, 0x00}},
		{1, RGB{0xff, 0xff, 0xff}},
		{2, RGB{0x89, 0x40, 0x36}},
		{3, RGB{0x7a, 0xbf, 0xc7}},
		{4, RGB{0x8a, 0x46, 0xae}},
		{5, RGB{0x68, 0xa9, 0x41}},
		{6, RGB{0x3e, 0x31, 0xa2}},
		{7, RGB{0xd0, 0xdc, 0x71}},
		{8, RGB{0x90, 0x5f, 0x25}},
		{9, RGB{0x5c, 0x47, 0x00}},
		{10, RGB{0xbb, 0x77, 0x6d}},
		{11, RGB{0x55, 0x55, 0x55}},
		{12, RGB{0x80, 0x80, 0x80}},
		{13, RGB{0xac, 0xea, 0x88}},
		{14, RGB{0x7c, 0x70, 0xda}},
		{15, RGB{0xab, 0xab, 0xab}},
	},
	"electric cocillana": {
		{0, RGB{0x00, 0x00, 0x00}},
		{1, RGB{0xff, 0xff, 0xff}},
		{2, RGB{0x8b, 0x1f, 0x00}},
		{3, RGB{0x6f, 0xdf, 0xb7}},
		{4, RGB{0xa7, 0x3b, 0x9f}},
		{5, RGB{0x4a, 0xb5, 0x10}},
		{6, RGB{0x08, 0x00, 0x94}},
		{7, RGB{0xf3, 0xeb, 0x5b}},
		{8, RGB{0xa5, 0x42, 0x00}},
		{9, RGB{0x63, 0x29, 0x18}},
		{10, RGB{0xcb, 0x7b, 0x6f}},
		{11, RGB{0x45, 0x44, 0x44}},
		{12, RGB{0x9f, 0x9f, 0x9f}},
		{13, RGB{0x94, 0xff, 0x94}},
		{14, RGB{0x4a, 0x94, 0xd6}},
		{15, RGB{0xbd, 0xbd, 0xbd}},
	},
	"ste": {
		{0, RGB{0x00, 0x00, 0x00}},
		{1, RGB{0xff, 0xff, 0xff}},
		{2, RGB{0xc8, 0x35, 0x35}},
		{3, RGB{0x83, 0xf0, 0xdc}},
		{4, RGB{0xcc, 0x59, 0xc6}},
		{5, RGB{0x59, 0xcd, 0x36}},
		{6, RGB{0x41, 0x37, 0xcd}},
		{7, RGB{0xf7, 0xee, 0x59}},
		{8, RGB{0xd1, 0x7f, 0x30}},
		{9, RGB{0x91, 0x5f, 0x33}},
		{10, RGB{0xf9, 0x9b, 0x97}},
		{11, RGB{0x5b, 0x5b, 0x5b}},
		{12, RGB{0x8e, 0x8e, 0x8e}},
		{13, RGB{0x9d, 0xff, 0x9d}},
		{14, RGB{0x75, 0xa1, 0xec}},
		{15, RGB{0xc1, 0xc1, 0xc1}},
	},
	"perplex 1": {
		{0, RGB{0x00, 0x00, 0x00}},
		{1, RGB{0xfd, 0xfe, 0xfc}},
		{2, RGB{0xbe, 0x1a, 0x24}},
		{3, RGB{0x30, 0xe6, 0xc6}},
		{4, RGB{0xb4, 0x1a, 0xe2}},
		{5, RGB{0x1f, 0xd2, 0x1e}},
		{6, RGB{0x21, 0x1b, 0xae}},
		{7, RGB{0xdf, 0xf6, 0x0a}},
		{8, RGB{0xb8, 0x41, 0x04}},
		{9, RGB{0x6a, 0x33, 0x04}},
		{10, RGB{0xfe, 0x4a, 0x57}},
		{11, RGB{0x42, 0x45, 0x40}},
		{12, RGB{0x70, 0x74, 0x6f}},
		{13, RGB{0x59, 0xfe, 0x59}},
		{14, RGB{0x5f, 0x53, 0xfe}},
		{15, RGB{0xa4, 0xa7, 0xa2}},
	},
	"perplex 2": {
		{0, RGB{0x00, 0x00, 0x00}},
		{1, RGB{0xff, 0xff, 0xff}},
		{2, RGB{0xcd, 0x31, 0x00}},
		{3, RGB{0x81, 0xff, 0xd8}},
		{4, RGB{0xeb, 0x4c, 0xe1}},
		{5, RGB{0x69, 0xf7, 0x00}},
		{6, RGB{0x2a, 0x19, 0xdd}},
		{7, RGB{0xff, 0xff, 0x5b}},
		{8, RGB{0xe4, 0x6a, 0x00}},
		{9, RGB{0x75, 0x48, 0x2a}},
		{10, RGB{0xff, 0xa1, 0x96}},
		{11, RGB{0x70, 0x70, 0x70}},
		{12, RGB{0xb0, 0xb0, 0xb0}},
		{13, RGB{0xb9, 0xff, 0xb9}},
		{14, RGB{0x69, 0xc4, 0xff}},
		{15, RGB{0xeb, 0xeb, 0xeb}},
	},
}
