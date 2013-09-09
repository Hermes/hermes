/* paq9a archiver, Dec. 31, 2007 (C) 2007, Matt Mahoney

    LICENSE

    This program is free software; you can redistribute it and/or
    modify it under the terms of the GNU General Public License as
    published by the Free Software Foundation; either version 3 of
    the License, or (at your option) any later version.

    This program is distributed in the hope that it will be useful, but
    WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
    General Public License for more details at
    Visit <http://www.gnu.org/copyleft/gpl.html>.

paq9a is an experimental file compressor and archiver.  Usage:

  paq9a {a|x|l} archive [[-opt] files...]...

Commands:

  a = create archive and compress named files.
  x = extract from archive.
  l = list contents.

Archives are "solid".  You can only create new archives.  You cannot
modify existing archives.  File names are stored and extracted exactly as
named when the archive is created, but you have the option to rename them
during extraction.  Files are never clobbered.

The "a" command creates a new archive and adds the named files.
Wildcards are permitted if compiled with g++.  Options
and filenames may be in any order.  Options apply only to filenames
after the option, and override previous options.  Options are:

  -s = store without compression.
  -c = compress (default).
  -1 through -9 selects memory level from 18 MB to 1.5 GB  Default is -7
     using 405 MB.  The memory option must be set before the first file.
     Decompression requires the same amount of memory.

For example:

  paq9a a foo.paq9a a.txt -3 -s b.txt -c c.txt tmp/d.txt /tmp/e.txt

creates the archive foo.paq9a with 5 files.  The file b.txt is
stored without compression.  The other 4 files are compressed
at memory level 3.  Extraction requires the same memory as compression.

If any named file does not exist, then it is omitted from the archive
with a warning and the remaining files are added.  An existing
archive cannot be overwritten.  There must be at least one filename on
the command line.

The "x" command extracts the archive contents, creating files exactly
as named when the archive was created.  Files cannot be overwritten.
If a file already exists or cannot be created, then it is skipped.
For example, "tmp/d.txt" would be skipped if either the current
directory does not have a subdirectory tmp, or tmp is write
protected, or tmp/d.txt already exists.

If "x" is followed by one or more file names, then the output files
are renamed in the order they were added to the archive and any remaining
contents are extracted without renaming.  For example:

  paq9a x foo.paq9a x.txt y.txt

would extract a.txt to x.txt and b.txt to y.txt, then extract c.txt, 
tmp/d.txt and /tmp/e.txt.  If the command line has more filenames than
the archive then the extra arguments are ignored.  Options are not
allowed.

The "l" (letter l) command lists the contents.  Any extra arguments
are ignored.

Any other command, or no command, displays a help message.


ARCHIVE FORMAT

  "lPq" 1 mem [filename {'\0' mode usize csize contents}...]...

The first 4 bytes are "lPq\x01" (1 is the version number).

mem is a digit '1' through '9', where '9' uses the most memory (1.5 GB).

A file is stored as one or more blocks.  The filename is stored
only in the first block as a NUL terminated string.  Subsequent
blocks start with a 0.

The mode is 's' if the block is stored and 'c' if compressed.

usize = uncompressed size as a 4 byte big-endian number (MSB first).

csize = compressed size as a 4 byte big-endian number.

The contents is copied from the file itself if mode is 's' or the
compressed contents otherwise.  Its length is exactly csize bytes.


COMPRESSED FORMAT

Files are preprocessed with LZP and then compressed with a context
mixing compressor and arithmetic coded one bit at a time.  Model
contents are maintained across files.

The LZP stage predicts the next byte by matching the current context
(order 12 or higher) to a rotating buffer.  If a match is found
then the next byte after the match is predicted.  If the next byte
matches the prediction, then a 1 bit is coded and the context is extended.
Otherwise a 0 is coded followed by 8 bits of the actual byte in MSB to 
LSB order.

A 1 bit is modeled using the match length as context, then refined
in 3 stages using sucessively longer contexts.  The predictions are 
adjusted by 2 input neurons selected by a context hash with the second 
input fixed.

If the LZP prediction is missed, then the literal is coded using a chain
of predicions which are mixed using neurons, where one input is the
previous prediction and the second input is the prediction given the
current context.  The current context is mapped to an 8 bit state
representing the bit history, the sequence of bits previously observed
in that context.  The bit history is used both to select the neuron
and is mapped to a prediction that provides the second input.  In addition,
if the known bits of the current byte match the LZP incorrectly predicted
byte, then this fact is used to select one of 2 sets of neurons (512 total).

The contexts, in order, are sparse order-1 with gaps of 3, 2, and 1
byte, then orders 1 through 6, then word orders 0 and 1, where a word
is a sequenece of case insensitive letters (useful for compressing text).
Contexts longer than 1 are hashed.  Order-n contexts consist of a hash
of the last n bytes plus the 0 to 7 known bits of the current byte.
The order 6 context and the word order 0 and 1 contexts also include
the LZP predicted byte.

All mixing is in the logistic or "stretched" domain: stretch(p) = ln(p/(1-p)),
then "squashed" by the inverse function: squash(p) = 1/(1 + exp(-p)) before
arithmetic coding.  A 2 input neuron has 2 weights (w0 and w1)
selected by context.  Given inputs x0 and x1 (2 predictions, or one
prediction and a constant), the output prediction is computed:
p = w0*x0 + w1*x1.  If the actual bit is y, then the weights are updated
to minimize its coding cost:

  error = y - squash(p)
  w0 += x0 * error * L
  w1 += x1 * error * L

where L is the learning rate, normally 1/256, but increased by a factor
of 4 an 2 for the first 2 training cycles (using the 2 low bits
of w0 as a counter).  In the implementation, p is represented by a fixed
point number with a 12 bit fractional part in the linear domain (0..4095)
and 8 bits in the logistic domain (-2047..2047 representing -8..8).
Weights are scaled by 24 bits.  Both weights are initialized to 1/2,
expecting 2 probabilities, weighted equally).  However, when one input
(x0) is fixed, its weight (w0) is initialized to 0.

A bit history represents the sequence of 0 and 1 bits observed in a given
context.  An 8 bit state represents all possible sequences up to 4 bits
long.  Longer sequences are represented by a count of 0 and 1 bits, plus
an indicator of the most recent bit.  If counts grow too large, then the
next state represents a pair of smaller counts with about the same ratio.
The state table is the same as used in PAQ8 (all versions) and LPAQ1.

A state is mapped to a prediction by using a table.  A table entry
contains 2 values, p, initialized to 1/2, and n, initialized to 0.
The output prediciton is p (in the linear domain, not stretched).
If the actual bit is y, then the entry is updated:

  error = y - p
  p += error/(n + 1.5)
  if n < limit then n += 1

In practice, p is scaled by 22 bits, and n is 10 bits, packed into
one 32 bit integer.  The limit is 255.

Every 4 bits, contexts are mapped to arrays of 15 states using a 
hash table.  The first element is the bit history for the current
context ending on a half byte boundary, followed by all possible
contexts formed by appending up to 3 more bits.

A hash table accepts a 32 bit context, which must be a hash if
longer than 4 bytes.  The input is further hashed and divided into
an index (depending on the table size, a power of 2), and an 8 bit
checksum which is stored in the table and used to detect collisions
(not perfectly).  A lookup tests 3 adjacent locations within a single
64 byte cache line, and if a matching checksum is not found, then the
entry with the smallest value in the first data element is replaced
(LFU replacement policy).  This element represents a bit history
for a context ending on a half byte boundary.  The states are ordered
so that larger values represent larger total bit counts, which
estimates the likelihood of future use.  The initial state is 0.

Memory is allocated from MEM = pow(2, opt+22) bytes, where opt is 1 through
9 (user selected).  Of this, MEM/2 is for the hash table for storing literal
context states, MEM/8 for the rotating LZP buffer, and MEM/8 for a 
hash table of pointers into the buffer, plus 12 MB for miscellaneous data.
Total memory usage is 0.75*MEM + 12 MB.


ARITHMETIC CODING

The arithmetic coder codes a bit with probability p using log2(1/p) bits.
Given input string y, the output is a binary fraction x such that
P(< y) <= x < P(<= y) where P(< y) means the total probability of all inputs
lexicographically less than y and P(<= y) = P(< y) + P(y).  Note that one
can always find x with length at most log2(P(y)) + 1 bits.

x can be computed efficiently by maintaining a range, low <= x < high
(initially 0..1) and expressing P(y) as a product of predictions:
P(y) = P(y1) P(y2|y1) P(y3|y1y2) P(y4|y1y2y3) ... P(yn|y1y2...yn-1)
where the term P(yi|y0y1...yi-1) means the probability that yi is 1
given the context y1...yi-1, the previous i-1 bits of y.  For each
prediction p, the range is split in proportion to the probabilities
of 0 and 1, then updated by taking the half corresponding to the actual
bit y as the new range, i.e.

  mid = low + (high - low) * p(y = 1)
  if y = 0 then (low, high) := (mid, high)
  if y = 1 then (low, high) := (low, mid)

As low and high approach each other, the high order bits of x become
known (because they are the same throughout the range) and can be
output immediately.

For decoding, the range is split as before and the range is updated
to the half containing x.  The corresponding bit y is used to update
the model.  Thus, the model has the same knowledge for coding and
decoding.

*/

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <ctype.h>
#define NDEBUG  // remove for debugging
#include <assert.h>

int allocated=0;  // Total memory allocated by alloc()

// Create an array p of n elements of type T
template <class T> void alloc(T*&p, int n) {
  p=(T*)calloc(n, sizeof(T));
  if (!p) printf("Out of memory\n"), exit(1);
  allocated+=n*sizeof(T);
}

// 8, 16, 32 bit unsigned types (adjust as appropriate)
typedef unsigned char  U8;
typedef unsigned short U16;
typedef unsigned int   U32;

///////////////////////////// Squash //////////////////////////////

// return p = 1/(1 + exp(-d)), d scaled by 8 bits, p scaled by 12 bits
class Squash {
  short tab[4096];
public:
  Squash();
  int operator()(int d) {
    d+=2048;
    if (d<0) return 0;
    else if (d>4095) return 4095;
    else return tab[d];
  }
} squash;

Squash::Squash() {
  static const int t[33]={
    1,2,3,6,10,16,27,45,73,120,194,310,488,747,1101,
    1546,2047,2549,2994,3348,3607,3785,3901,3975,4022,
    4050,4068,4079,4085,4089,4092,4093,4094};
  for (int i=-2048; i<2048; ++i) {
    int w=i&127;
    int d=(i>>7)+16;
    tab[i+2048]=(t[d]*(128-w)+t[(d+1)]*w+64) >> 7;
  }
}

//////////////////////////// Stretch ///////////////////////////////

// Inverse of squash. stretch(d) returns ln(p/(1-p)), d scaled by 8 bits,
// p by 12 bits.  d has range -2047 to 2047 representing -8 to 8.  
// p has range 0 to 4095 representing 0 to 1.

class Stretch {
  short t[4096];
public:
  Stretch();
  int operator()(int p) const {
    assert(p>=0 && p<4096);
    return t[p];
  }
} stretch;

Stretch::Stretch() {
  int pi=0;
  for (int x=-2047; x<=2047; ++x) {  // invert squash()
    int i=squash(x);
    for (int j=pi; j<=i; ++j)
      t[j]=x;
    pi=i+1;
  }
  t[4095]=2047;
}

///////////////////////////// ilog //////////////////////////////

// ilog(x) = round(log2(x) * 16), 0 <= x < 64K
class Ilog {
  U8* t;
public:
  int operator()(U16 x) const {return t[x];}
  Ilog();
} ilog;

// Compute lookup table by numerical integration of 1/x
Ilog::Ilog() {
  alloc(t, 65536);
  U32 x=14155776;
  for (int i=2; i<65536; ++i) {
    x+=774541002/(i*2-1);  // numerator is 2^29/ln 2
    t[i]=x>>24;
  }
}

// llog(x) accepts 32 bits
inline int llog(U32 x) {
  if (x>=0x1000000)
    return 256+ilog(x>>16);
  else if (x>=0x10000)
    return 128+ilog(x>>8);
  else
    return ilog(x);
}

///////////////////////// state table ////////////////////////

// State table:
//   nex(state, 0) = next state if bit y is 0, 0 <= state < 256
//   nex(state, 1) = next state if bit y is 1
//
// States represent a bit history within some context.
// State 0 is the starting state (no bits seen).
// States 1-30 represent all possible sequences of 1-4 bits.
// States 31-252 represent a pair of counts, (n0,n1), the number
//   of 0 and 1 bits respectively.  If n0+n1 < 16 then there are
//   two states for each pair, depending on if a 0 or 1 was the last
//   bit seen.
// If n0 and n1 are too large, then there is no state to represent this
// pair, so another state with about the same ratio of n0/n1 is substituted.
// Also, when a bit is observed and the count of the opposite bit is large,
// then part of this count is discarded to favor newer data over old.

static const U8 State_table[256][2]={
{  1,  2},{  3,  5},{  4,  6},{  7, 10},{  8, 12},{  9, 13},{ 11, 14}, // 0
{ 15, 19},{ 16, 23},{ 17, 24},{ 18, 25},{ 20, 27},{ 21, 28},{ 22, 29}, // 7
{ 26, 30},{ 31, 33},{ 32, 35},{ 32, 35},{ 32, 35},{ 32, 35},{ 34, 37}, // 14
{ 34, 37},{ 34, 37},{ 34, 37},{ 34, 37},{ 34, 37},{ 36, 39},{ 36, 39}, // 21
{ 36, 39},{ 36, 39},{ 38, 40},{ 41, 43},{ 42, 45},{ 42, 45},{ 44, 47}, // 28
{ 44, 47},{ 46, 49},{ 46, 49},{ 48, 51},{ 48, 51},{ 50, 52},{ 53, 43}, // 35
{ 54, 57},{ 54, 57},{ 56, 59},{ 56, 59},{ 58, 61},{ 58, 61},{ 60, 63}, // 42
{ 60, 63},{ 62, 65},{ 62, 65},{ 50, 66},{ 67, 55},{ 68, 57},{ 68, 57}, // 49
{ 70, 73},{ 70, 73},{ 72, 75},{ 72, 75},{ 74, 77},{ 74, 77},{ 76, 79}, // 56
{ 76, 79},{ 62, 81},{ 62, 81},{ 64, 82},{ 83, 69},{ 84, 71},{ 84, 71}, // 63
{ 86, 73},{ 86, 73},{ 44, 59},{ 44, 59},{ 58, 61},{ 58, 61},{ 60, 49}, // 70
{ 60, 49},{ 76, 89},{ 76, 89},{ 78, 91},{ 78, 91},{ 80, 92},{ 93, 69}, // 77
{ 94, 87},{ 94, 87},{ 96, 45},{ 96, 45},{ 48, 99},{ 48, 99},{ 88,101}, // 84
{ 88,101},{ 80,102},{103, 69},{104, 87},{104, 87},{106, 57},{106, 57}, // 91
{ 62,109},{ 62,109},{ 88,111},{ 88,111},{ 80,112},{113, 85},{114, 87}, // 98
{114, 87},{116, 57},{116, 57},{ 62,119},{ 62,119},{ 88,121},{ 88,121}, // 105
{ 90,122},{123, 85},{124, 97},{124, 97},{126, 57},{126, 57},{ 62,129}, // 112
{ 62,129},{ 98,131},{ 98,131},{ 90,132},{133, 85},{134, 97},{134, 97}, // 119
{136, 57},{136, 57},{ 62,139},{ 62,139},{ 98,141},{ 98,141},{ 90,142}, // 126
{143, 95},{144, 97},{144, 97},{ 68, 57},{ 68, 57},{ 62, 81},{ 62, 81}, // 133
{ 98,147},{ 98,147},{100,148},{149, 95},{150,107},{150,107},{108,151}, // 140
{108,151},{100,152},{153, 95},{154,107},{108,155},{100,156},{157, 95}, // 147
{158,107},{108,159},{100,160},{161,105},{162,107},{108,163},{110,164}, // 154
{165,105},{166,117},{118,167},{110,168},{169,105},{170,117},{118,171}, // 161
{110,172},{173,105},{174,117},{118,175},{110,176},{177,105},{178,117}, // 168
{118,179},{110,180},{181,115},{182,117},{118,183},{120,184},{185,115}, // 175
{186,127},{128,187},{120,188},{189,115},{190,127},{128,191},{120,192}, // 182
{193,115},{194,127},{128,195},{120,196},{197,115},{198,127},{128,199}, // 189
{120,200},{201,115},{202,127},{128,203},{120,204},{205,115},{206,127}, // 196
{128,207},{120,208},{209,125},{210,127},{128,211},{130,212},{213,125}, // 203
{214,137},{138,215},{130,216},{217,125},{218,137},{138,219},{130,220}, // 210
{221,125},{222,137},{138,223},{130,224},{225,125},{226,137},{138,227}, // 217
{130,228},{229,125},{230,137},{138,231},{130,232},{233,125},{234,137}, // 224
{138,235},{130,236},{237,125},{238,137},{138,239},{130,240},{241,125}, // 231
{242,137},{138,243},{130,244},{245,135},{246,137},{138,247},{140,248}, // 238
{249,135},{250, 69},{ 80,251},{140,252},{249,135},{250, 69},{ 80,251}, // 245
{140,252},{  0,  0},{  0,  0},{  0,  0}};  // 252
#define nex(state,sel) State_table[state][sel]

//////////////////////////// StateMap //////////////////////////

// A StateMap maps a context to a probability.  Methods:
//
// Statemap sm(n) creates a StateMap with n contexts using 4*n bytes memory.
// sm.p(cx, limit) converts state cx (0..n-1) to a probability (0..4095)
//     that the next updated bit y=1.
//     limit (1..1023, default 255) is the maximum count for computing a
//     prediction.  Larger values are better for stationary sources.
// sm.update(y) updates the model with actual bit y (0..1).

class StateMap {
protected:
  const int N;  // Number of contexts
  int cxt;      // Context of last prediction
  U32 *t;       // cxt -> prediction in high 22 bits, count in low 10 bits
  static int dt[1024];  // i -> 16K/(i+3)
public:
  StateMap(int n=256);

  // update bit y (0..1)
  void update(int y, int limit=255) {
    assert(cxt>=0 && cxt<N);
    int n=t[cxt]&1023, p=t[cxt]>>10;  // count, prediction
    if (n<limit) ++t[cxt];
    else t[cxt]=t[cxt]&0xfffffc00|limit;
    t[cxt]+=(((y<<22)-p)>>3)*dt[n]&0xfffffc00;
  }

  // predict next bit in context cx
  int p(int cx) {
    assert(cx>=0 && cx<N);
    return t[cxt=cx]>>20;
  }
};

int StateMap::dt[1024]={0};

StateMap::StateMap(int n): N(n), cxt(0) {
  alloc(t, N);
  for (int i=0; i<N; ++i)
    t[i]=1<<31;
  if (dt[0]==0)
    for (int i=0; i<1024; ++i)
      dt[i]=16384/(i+i+3);
}

//////////////////////////// Mix, APM /////////////////////////

// Mix combines 2 predictions and a context to produce a new prediction.
// Methods:
// Mix m(n) -- creates allowing with n contexts.
// m.pp(p1, p2, cx) -- inputs 2 stretched predictions and a context cx
//   (0..n-1) and returns a stretched prediction.  Stretched predictions
//   are fixed point numbers with an 8 bit fraction, normally -2047..2047
//   representing -8..8, such that 1/(1+exp(-p) is the probability that
//   the next update will be 1.
// m.update(y) updates the model after a prediction with bit y (0..1).

class Mix {
protected:
  const int N;  // n
  int* wt;  // weights, scaled 24 bits
  int x1, x2;    // inputs, scaled 8 bits (-2047 to 2047)
  int cxt;  // last context (0..n-1)
  int pr;   // last output
public:
  Mix(int n=512);
  int pp(int p1, int p2, int cx) {
    assert(cx>=0 && cx<N);
    cxt=cx*2;
    return pr=(x1=p1)*(wt[cxt]>>16)+(x2=p2)*(wt[cxt+1]>>16)+128>>8;
  }
  void update(int y) {
    assert(y==0 || y==1);
    int err=((y<<12)-squash(pr));
    if ((wt[cxt]&3)<3)
      err*=4-(++wt[cxt]&3);
    err=err+8>>4;
    wt[cxt]+=x1*err&-4;
    wt[cxt+1]+=x2*err;
  }
};

Mix::Mix(int n): N(n), x1(0), x2(0), cxt(0), pr(0) {
  alloc(wt, n*2);
  for (int i=0; i<N*2; ++i)
    wt[i]=1<<23;
}

// An APM is a Mix optimized for a constant in place of p1, used to
// refine a stretched prediction given a context cx. 
// Normally p1 is in the range (0..4095) and p2 is doubled.

class APM: public Mix {
public:
  APM(int n);
};

APM::APM(int n): Mix(n) {
  for (int i=0; i<n; ++i)
    wt[2*i]=0;
}

//////////////////////////// HashTable /////////////////////////

// A HashTable maps a 32-bit index to an array of B bytes.
// The first byte is a checksum using the upper 8 bits of the
// index.  The second byte is a priority (0 = empty) for hash
// replacement.  The index need not be a hash.

// HashTable<B> h(n) - create using n bytes  n and B must be 
//     powers of 2 with n >= B*4, and B >= 2.
// h[i] returns array [1..B-1] of bytes indexed by i, creating and
//     replacing another element if needed.  Element 0 is the
//     checksum and should not be modified.

template <int B>
class HashTable {
  U8* t;  // table: 1 element = B bytes: checksum priority data data
  const U32 N;  // size in bytes
public:
  HashTable(int n);
  ~HashTable();
  U8* operator[](U32 i);
};

template <int B>
HashTable<B>::HashTable(int n): t(0), N(n) {
  assert(B>=2 && (B&B-1)==0);
  assert(N>=B*4 && (N&N-1)==0);
  alloc(t, N+B*4+64);
  t+=64-int(((long)t)&63);  // align on cache line boundary
}

template <int B>
inline U8* HashTable<B>::operator[](U32 i) {
  i*=123456791;
  i=i<<16|i>>16;
  i*=234567891;
  int chk=i>>24;
  i=i*B&N-B;
  if (t[i]==chk) return t+i;
  if (t[i^B]==chk) return t+(i^B);
  if (t[i^B*2]==chk) return t+(i^B*2);
  if (t[i+1]>t[i+1^B] || t[i+1]>t[i+1^B*2]) i^=B;
  if (t[i+1]>t[i+1^B^B*2]) i^=B^B*2;
  memset(t+i, 0, B);
  t[i]=chk;
  return t+i;
}

template <int B>
HashTable<B>::~HashTable() {
  int c=0, c0=0;
  for (U32 i=0; i<N; ++i) {
    if (t[i]) {
      ++c;
      if (i%B==0) ++c0;
    }
  }
  printf("HashTable<%d> %1.4f%% full, %1.4f%% utilized of %d KiB\n",
    B, 100.0*c0*B/N, 100.0*c/N, N>>10);
}

////////////////////////// LZP /////////////////////////

U32 MEM=1<<29;  // Global memory limit, 1 << 22+(memory option)

// LZP predicts the next byte and maintains context.  Methods:
// c() returns the predicted byte for the next update, or -1 if none.
// p() returns the 12 bit probability (0..4095) that c() is next.
// update(ch) updates the model with actual byte ch (0..255).
// c(i) returns the i'th prior byte of context, i > 0.
// c4() returns the order 4 context, shifted into the LSB.
// c8() returns a hash of the order 8 context, shifted 4 bits into LSB.
// word0, word1 are hashes of the current and previous word (a-z).

class LZP {
private:
  const int N, H; // buf, t sizes
  enum {MINLEN=12};  // minimum match length
  U8* buf;     // Rotating buffer of size N
  U32* t;      // hash table of pointers in high 24 bits, state in low 8 bits
  int match;   // start of match
  int len;     // length of match
  int pos;     // position of next ch to write to buf
  U32 h;       // context hash
  U32 h1;      // hash of last 8 byte updates, shifting 4 bits to MSB
  U32 h2;      // last 4 updates, shifting 8 bits to MSB
  StateMap sm1; // len+offset -> p
  APM a1, a2, a3;   // p, context -> p
  int literals, matches;  // statistics
public:
  U32 word0, word1;  // hashes of last 2 words (case insensitive a-z)
  LZP();
  ~LZP();
  int c();     // predicted char
  int c(int i);// context
  int c4() {return h2;}  // order 4 context, c(1) in LSB
  int c8() {return h1;}  // hashed order 8 context
  int p();     // probability that next char is c() * 4096
  void update(int ch);  // update model with actual char ch
};

// Initialize
LZP::LZP(): N(MEM/8), H(MEM/32),
    match(-1), len(0), pos(0), h(0), h1(0), h2(0), 
    sm1(0x200), a1(0x10000), a2(0x40000), a3(0x100000),
    literals(0), matches(0), word0(0), word1(0) {
  assert(MEM>0);
  assert(H>0);
  alloc(buf, N);
  alloc(t, H);
}

// Print statistics
LZP::~LZP() {
  int c=0;
  for (int i=0; i<H; ++i)
    c+=(t[i]!=0);
  printf("LZP hash table %1.4f%% full of %d KiB\n"
    "LZP buffer %1.4f%% full of %d KiB\n", 
    100.0*c/H, H>>8, pos<N?100.0*pos/N:100.0, N>>10);
  printf("LZP %d literals, %d matches (%1.4f%% matched)\n",
    literals, matches, 
    literals+matches>0?100.0*matches/(literals+matches):0.0);
}

// Predicted next byte, or -1 for no prediction
inline int LZP::c() {
  return len>=MINLEN ? buf[match&N-1] : -1;
}

// Return i'th byte of context (i > 0)
inline int LZP::c(int i) {
  assert(i>0);
  return buf[pos-i&N-1];
}

// Return prediction that c() will be the next byte (0..4095)
int LZP::p() {
  if (len<MINLEN) return 0;
  int cxt=len;
  if (len>28) cxt=28+(len>=32)+(len>=64)+(len>=128);
  int pc=c();
  int pr=sm1.p(cxt);
  pr=stretch(pr);
  pr=a1.pp(2048, pr*2, h2*256+pc&0xffff)*3+pr>>2;
  pr=a2.pp(2048, pr*2, h1*(11<<6)+pc&0x3ffff)*3+pr>>2;
  pr=a3.pp(2048, pr*2, h1*(7<<4)+pc&0xfffff)*3+pr>>2;
  pr=squash(pr);
  return pr;
}

// Update model with predicted byte ch (0..255)
void LZP::update(int ch) {
  int y=c()==ch;     // 1 if prediction of ch was right, else 0
  h1=h1*(3<<4)+ch+1; // update context hashes
  h2=h2<<8|ch;
  h=h*(5<<2)+ch+1&H-1;
  if (len>=MINLEN) {
    sm1.update(y);
    a1.update(y);
    a2.update(y);
    a3.update(y);
  }
  if (isalpha(ch))
    word0=word0*(29<<2)+tolower(ch);
  else if (word0)
    word1=word0, word0=0;
  buf[pos&N-1]=ch;   // update buf
  ++pos;
  if (y) {  // extend match
    ++len;
    ++match;
    ++matches;
  }
  else {  // find new match, try order 6 context first
    ++literals;
    y=0;
    len=1;
    match=t[h];
    if (!((match^pos)&N-1)) --match;
    while (len<=128 && buf[match-len&N-1]==buf[pos-len&N-1]) ++len;
    --len;
  }
  t[h]=pos;
}

LZP* lzp=0;

//////////////////////////// Predictor /////////////////////////

// A Predictor estimates the probability that the next bit of
// uncompressed data is 1.  Methods:
// Predictor() creates.
// p() returns P(1) as a 12 bit number (0-4095).
// update(y) trains the predictor with the actual bit (0 or 1).

class Predictor {
  enum {N=11}; // number of contexts
  int c0;      // last 0-7 bits with leading 1, 0 before LZP flag
  int nibble;  // last 0-3 bits with leading 1 (1..15)
  int bcount;  // number of bits in c0 (0..7)
  HashTable<16> t;  // context -> state
  StateMap sm[N];   // state -> prediction
  U8* cp[N];   // i -> state array of bit histories for i'th context
  U8* sp[N];   // i -> pointer to bit history for i'th context
  Mix m[N-1];  // combines 2 predictions given a context
  APM a1, a2, a3;  // adjusts a prediction given a context
  U8* t2;      // order 1 contexts -> state

public:
  Predictor();
  int p();
  void update(int y);
};

// Initialize
Predictor::Predictor():
    c0(0), nibble(1), bcount(0), t(MEM/2),
    a1(0x10000), a2(0x10000), a3(0x10000) {
  alloc(t2, 0x40000);
  for (int i=0; i<N; ++i)
    sp[i]=cp[i]=t2;
}

// Update model
void Predictor::update(int y) {
  assert(y==0 || y==1);
  assert(bcount>=0 && bcount<8);
  assert(c0>=0 && c0<256);
  assert(nibble>=1 && nibble<=15);
  if (c0==0)
    c0=1-y;
  else {
    *sp[0]=nex(*sp[0], y);
    sm[0].update(y);
    for (int i=1; i<N; ++i) {
      *sp[i]=nex(*sp[i], y);
      sm[i].update(y);
      m[i-1].update(y);
    }
    c0+=c0+y;
    if (++bcount==8) bcount=c0=0;
    if ((nibble+=nibble+y)>=16) nibble=1;
    a1.update(y);
    a2.update(y);
    a3.update(y);
  }
}

// Predict next bit
int Predictor::p() {
  assert(lzp);
  if (c0==0)
    return lzp->p();
  else {

    // Set context pointers
    int pc=lzp->c();  // mispredicted byte
    int r=pc+256>>8-bcount==c0;  // c0 consistent with mispredicted byte?
    U32 c4=lzp->c4();  // last 4 whole context bytes, shifted into LSB
    U32 c8=(lzp->c8()<<4)-1;  // hash of last 7 bytes with 4 trailing 1 bits
    if ((bcount&3)==0) {  // nibble boundary?  Update context pointers
      pc&=-r;
      U32 c4p=c4<<8;
      if (bcount==0) {  // byte boundary?  Update order-1 context pointers
        cp[0]=t2+(c4>>16&0xff00);
        cp[1]=t2+(c4>>8 &0xff00)+0x10000;
        cp[2]=t2+(c4    &0xff00)+0x20000;
        cp[3]=t2+(c4<<8 &0xff00)+0x30000;
      }
      cp[4]=t[(c4p&0xffff00)-c0];
      cp[5]=t[(c4p&0xffffff00)*3+c0];
      cp[6]=t[c4*7+c0];
      cp[7]=t[(c8*5&0xfffffc)+c0];
      cp[8]=t[(c8*11&0xffffff0)+c0+pc*13];
      cp[9]=t[lzp->word0*5+c0+pc*17];
      cp[10]=t[lzp->word1*7+lzp->word0*11+c0+pc*37];
    }

    // Mix predictions
    r<<=8;
    sp[0]=&cp[0][c0];
    int pr=stretch(sm[0].p(*sp[0]));
    for (int i=1; i<N; ++i) {
      sp[i]=&cp[i][i<4?c0:nibble];
      int st=*sp[i];
      pr=m[i-1].pp(pr, stretch(sm[i].p(st)), st+r)*3+pr>>2;
    }
    pr=a1.pp(512, pr*2, c0+pc*256&0xffff)*3+pr>>2;  // Adjust prediction
    pr=a2.pp(512, pr*2, c4<<8&0xff00|c0)*3+pr>>2;
    pr=a3.pp(512, pr*2, c4*3+c0&0xffff)*3+pr>>2;
    return squash(pr);
  }
}

Predictor* predictor=0;

/////////////////////////// get4, put4 //////////////////////////

// Read/write a 4 byte big-endian number
int get4(FILE* in) {
  int r=getc(in);
  r=r*256+getc(in);
  r=r*256+getc(in);
  r=r*256+getc(in);
  return r;
}

void put4(U32 c, FILE* out) {
  fprintf(out, "%c%c%c%c", c>>24, c>>16, c>>8, c);
}

//////////////////////////// Encoder ////////////////////////////

// An Encoder arithmetic codes in blocks of size BUFSIZE.  Methods:
// Encoder(COMPRESS, f) creates encoder for compression to archive f, which
//     must be open past any header for writing in binary mode.
// Encoder(DECOMPRESS, f) creates encoder for decompression from archive f,
//     which must be open past any header for reading in binary mode.
// code(i) in COMPRESS mode compresses bit i (0 or 1) to file f.
// code() in DECOMPRESS mode returns the next decompressed bit from file f.
// count() should be called after each byte is compressed.
// flush() should be called after compression is done.  It is also called
//   automatically when a block is written.

typedef enum {COMPRESS, DECOMPRESS} Mode;
class Encoder {
private:
  const Mode mode;       // Compress or decompress?
  FILE* archive;         // Compressed data file
  U32 x1, x2;            // Range, initially [0, 1), scaled by 2^32
  U32 x;                 // Decompress mode: last 4 input bytes of archive
  enum {BUFSIZE=0x20000};
  static unsigned char* buf; // Compression output buffer, size BUFSIZE
  int usize, csize;      // Buffered uncompressed and compressed sizes
  double usum, csum;     // Total of usize, csize

public:
  Encoder(Mode m, FILE* f);
  void flush();  // call this when compression is finished

  // Compress bit y or return decompressed bit
  int code(int y=0) {
    assert(predictor);
    int p=predictor->p();
    assert(p>=0 && p<4096);
    p+=p<2048;
    U32 xmid=x1 + (x2-x1>>12)*p + ((x2-x1&0xfff)*p>>12);
    assert(xmid>=x1 && xmid<x2);
    if (mode==DECOMPRESS) y=x<=xmid;
    y ? (x2=xmid) : (x1=xmid+1);
    predictor->update(y);
    while (((x1^x2)&0xff000000)==0) {  // pass equal leading bytes of range
      if (mode==COMPRESS) buf[csize++]=x2>>24;
      x1<<=8;
      x2=(x2<<8)+255;
      if (mode==DECOMPRESS) x=(x<<8)+getc(archive);
    }
    return y;
  }

  // Count one byte
  void count() {
    assert(mode==COMPRESS);
    ++usize;
    if (csize>BUFSIZE-256)
      flush();
  }
};
unsigned char* Encoder::buf=0;

// Create in mode m (COMPRESS or DECOMPRESS) with f opened as the archive.
Encoder::Encoder(Mode m, FILE* f):
    mode(m), archive(f), x1(0), x2(0xffffffff), x(0), 
    usize(0), csize(0), usum(0), csum(0) {
  if (mode==DECOMPRESS) {  // x = first 4 bytes of archive
    for (int i=0; i<4; ++i)
      x=(x<<8)+(getc(archive)&255);
    csize=4;
  }
  else if (!buf)
    alloc(buf, BUFSIZE);
}

// Write a compressed block and reinitialize the encoder.  The format is:
//   uncompressed size (usize, 4 byte, MSB first)
//   compressed size (csize, 4 bytes, MSB first)
//   compressed data (csize bytes)
void Encoder::flush() {
  if (mode==COMPRESS) {
    buf[csize++]=x1>>24;
    buf[csize++]=255;
    buf[csize++]=255;
    buf[csize++]=255;
    putc(0, archive);
    putc('c', archive);
    put4(usize, archive);
    put4(csize, archive);
    fwrite(buf, 1, csize, archive);
    usum+=usize;
    csum+=csize+10;
    printf("%15.0f -> %15.0f"
      "\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b", 
      usum, csum);
    x1=x=usize=csize=0;
    x2=0xffffffff;
  }
}

/////////////////////////// paq9a ////////////////////////////////

// Compress or decompress from in to out, depending on whether mode
// is COMPRESS or DECOMPRESS.  A byte c is encoded as a 1 bit if it
// is predicted by LZP, otherwise a 0 followed by 8 bits from MSB to LSB.
void paq9a(FILE* in, FILE* out, Mode mode) {
  if (!lzp && !predictor) {
    lzp=new LZP;
    predictor=new Predictor;
    printf("%8d KiB\b\b\b\b\b\b\b\b\b\b\b\b", allocated>>10);
  }
  if (mode==COMPRESS) {
    Encoder e(COMPRESS, out);
    int c;
    while ((c=getc(in))!=EOF) {
      int cp=lzp->c();
      if (c==cp)
        e.code(1);
      else
        for (int i=8; i>=0; --i)
          e.code(c>>i&1);
      e.count();
      lzp->update(c);
    }
    e.flush();
  }
  else {  // DECOMPRESS
    int usize=get4(in);
    get4(in);  // csize
    Encoder e(DECOMPRESS, in);
    while (usize--) {
      int c=lzp->c();
      if (e.code()==0) {
        c=1;
        while (c<256) c+=c+e.code();
        c&=255;
      }
      if (out) putc(c, out);
      lzp->update(c);
    }
  }
}


///////////////////////////// store ///////////////////////////

// Store a file in blocks as: {'\0' mode usize csize contents}...
void store(FILE* in, FILE* out) {
  assert(in);
  assert(out);

  // Store in blocks
  const int BLOCKSIZE=0x100000;
  static char* buf=0;
  if (!buf) alloc(buf, BLOCKSIZE);
  bool first=true;
  while (true) {
    int n=fread(buf, 1, BLOCKSIZE, in);
    if (!first && n<=0) break;
    fprintf(out, "%c%c", 0, 's');
    put4(n, out);  // usize
    put4(n, out);  // csize
    fwrite(buf, 1, n, out);
    first=false;
  }

  // Close file
  fclose(in);
}

// Write usize == csize bytes of an uncompressed block from in to out
void unstore(FILE* in, FILE* out) {
  assert(in);
  int usize=get4(in);
  int csize=get4(in);
  if (usize!=csize)
    printf("Bad archive format: usize=%d csize=%d\n", usize, csize);
  static char* buf=0;
  const int BUFSIZE=0x1000;
  if (!buf) alloc(buf, BUFSIZE);
  while (csize>0) {
    usize=csize;
    if (usize>BUFSIZE) usize=BUFSIZE;
    if (int(fread(buf, 1, usize, in))!=usize)
      printf("Unexpected end of archive\n"), exit(1);
    if (out) fwrite(buf, 1, usize, out);
    csize-=usize;
  }
}

//////////////////////// Archiving functions ////////////////////////

const int MAXNAMELEN=1023;  // max filename length

// Return true if the first 4 bytes of in are a valid archive
bool check_archive(FILE* in) {
  return getc(in)=='p' && getc(in)=='Q' && getc(in)=='9' && getc(in)==1;
}

// Open archive and check for valid archive header, exit if bad.
// Set MEM to memory option '1' through '9'
FILE* open_archive(const char* filename) {
  FILE* in=fopen(filename, "rb");
  if (!in)
    printf("Cannot find archive %s\n", filename), exit(1);
  if (!check_archive(in) || (MEM=getc(in))<'1' || MEM>'9') {
    fclose(in);
    printf("%s: Not a paq9a archive\n", filename);
    exit(1);
  }
  return in;
}

// Compress filename to out.  option is 'c' to compress or 's' to store.
void compress(const char* filename, FILE* out, int option) {

  // Open input file
  FILE* in=fopen(filename, "rb");
  if (!in) {
    printf("File not found: %s\n", filename);
    return;
  }
  fprintf(out, "%s", filename);
  printf("%-40s ", filename);

  // Compress depending on option
  if (option=='s')
    store(in, out);
  else if (option=='c')
    paq9a(in, out, COMPRESS);
  printf("\n");
}

// List archive contents
void list(const char* archive) {
  double usum=0, csum=0;  // uncompressed and compressed size per file
  double utotal=0, ctotal=4;  // total size in archive
  static char filename[MAXNAMELEN+1];
  int mode=0;

  FILE* in=open_archive(archive);
  printf("\npaq9a -%c\n", MEM);
  while (true) {

    // Get filename, mode
    int c=getc(in);
    if (c==EOF) break;
    if (c) {   // start of new file?  Print previous file
      if (mode)
        printf("%10.0f -> %10.0f %c %s\n", usum, csum, mode, filename);
      int len=0;
      filename[len++]=c;
      while ((c=getc(in))!=EOF && c)
        if (len<MAXNAMELEN) filename[len++]=c;
      filename[len]=0;
      utotal+=usum;
      ctotal+=csum;
      usum=0;
      csum=len;
    }

    // Get uncompressed size
    mode=getc(in);
    int usize=get4(in);
    usum+=usize;

    // Get compressed size
    int csize=get4(in);
    csum+=csize+10;

    if (usize<0 || csize<0 || mode!='c' && mode!='s')
      printf("Archive corrupted usize=%d csize=%d mode=%d at %ld\n",
        usize, csize, mode, ftell(in)), exit(1);

    // Skip csize bytes
    const int BUFSIZE=0x1000;
    char buf[BUFSIZE];
    while (csize>BUFSIZE)
      csize-=fread(buf, 1, BUFSIZE, in);
    fread(buf, 1, csize, in);
  }
  printf("%10.0f -> %10.0f %c %s\n", usum, csum, mode, filename);
  utotal+=usum;
  ctotal+=csum;
  printf("%10.0f -> %10.0f total\n", utotal, ctotal);
  fclose(in);
}

// Extract files given command line arguments
// Input format is: [filename {'\0' mode usize csize contents}...]...
void extract(int argc, char** argv) {
  assert(argc>2);
  assert(argv[1][0]=='x');
  static char filename[MAXNAMELEN+1];  // filename from archive

  // Open archive
  FILE* in=open_archive(argv[2]);
  MEM=1<<22+MEM-'0';

  // Extract files
  argc-=3;
  argv+=3;
  FILE* out=0;
  while (true) {  // for each block

    // Get filename
    int c;
    for (int i=0;; ++i) {
      c=getc(in);
      if (c==EOF) break;
      if (i<MAXNAMELEN) filename[i]=c;
      if (!c) break;
    }
    if (c==EOF) break;

    // Open output file
    if (filename[0]) {  // new file?
      const char* fn=filename;
      if (argc>0) fn=argv[0], --argc, ++argv;
      if (out) fclose(out);
      out=fopen(fn, "rb");
      if (out) {
        printf("\nCannot overwrite file, skipping: %s ", fn);
        fclose(out);
        out=0;
      }
      else {
        out=fopen(fn, "wb");
        if (!out) printf("\nCannot create file: %s ", fn);
      }
      if (out) {
        if (fn==filename) printf("\n%s ", filename);
        else printf("\n%s -> %s ", filename, fn);
      }
    }

    // Extract block
    int mode=getc(in);
    if (mode=='s')
      unstore(in, out);
    else if (mode=='c')
      paq9a(in, out, DECOMPRESS);
    else
      printf("\nUnsupported compression mode %c %d at %ld\n", 
        mode, mode, ftell(in)), exit(1);
  }
  printf("\n");
  if (out) fclose(out);
}

// Command line is: paq9a {a|x|l} archive [[-option] files...]...
int main(int argc, char** argv) {
  clock_t start=clock();

  // Check command line arguments
  if (argc<3 || argv[1][1] || (argv[1][0]!='a' && argv[1][0]!='x'
      && argv[1][0]!='l') || (argv[1][0]=='a' && argc<4) || argv[2][0]=='-')
  {
    printf("paq9a archiver (C) 2007, Matt Mahoney\n"
      "Free software under GPL, http://www.gnu.org/copyleft/gpl.html\n"
      "\n"
      "To create archive: paq9a a archive [-1..-9] [[-s|-c] files...]...\n"
      "  -1..-9 = use 18 to 1585 MiB memory (default -7 = 408 MiB)\n"
      "  -s = store, -c = compress (default)\n"
      "To extract files:  paq9a x archive [files...]\n"
      "To list contents:  paq9a l archive\n");
    exit(1);
  }

  // Create archive
  if (argv[1][0]=='a') {
    int option = 'c';  // -c or -s
    FILE* out=fopen(argv[2], "rb");
    if (out) printf("Cannot overwrite archive %s\n", argv[2]), exit(1);
    out=fopen(argv[2], "wb");
    if (!out) printf("Cannot create archive %s\n", argv[2]), exit(1);
    fprintf(out, "pQ9%c", 1);
    int i=3;
    if (argc>3 && argv[3][0]=='-' && argv[3][1]>='1' && argv[3][1]<='9'
        && argv[3][2]==0) {
      putc(argv[3][1], out);
      MEM=1<<22+argv[3][1]-'0';
      ++i;
    }
    else
      putc('7', out);
    for (; i<argc; ++i) {
      if (argv[i][0]=='-' && (argv[i][1]=='c' || argv[i][1]=='s')
          && argv[i][2]==0)
        option=argv[i][1];
      else
        compress(argv[i], out, option);
    }
    printf("-> %ld in %1.2f sec\n", ftell(out),
      double(clock()-start)/CLOCKS_PER_SEC);
  }

  // List archive contents
  else if (argv[1][0]=='l')
    list(argv[2]);

  // Extract from archive
  else if (argv[1][0]=='x') {
    extract(argc, argv);
    printf("%1.2f sec\n", double(clock()-start)/CLOCKS_PER_SEC);
  }

  // Report statistics
  delete predictor;
  delete lzp;
  printf("Used %d KiB memory\n", allocated>>10);
  return 0;
}
