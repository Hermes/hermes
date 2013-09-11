# basic fountain code implementation
from random import randint

class block:
	
	def __init__(self, p, c):
		self.parents = p
		self.content = c
	
def bin_ascii(s):
	return ''.join(chr(int(s[i:i+8], 2)) for i in xrange(0, len(s), 8))

def encode(s, d, b, p):
	# s~ource, d~istribution, b~lock size, p~ercentage of blocks

	# convert s to binary!
	
	chunks = []
	for i in range(len(s) / b):
		chunks.append(s[i * b:(i * b) + b])

	blocks = []
	for i in range(int(len(chunks) * p)):
		index = randint(0, len(chunks)-1)
		iblock = block([index], chunks[index])
		for i in range(randint(0, d)):
			index = randint(0, len(chunks)-1)
			iblock.parents.append(index)
			iblock.content = xor(iblock.content, chunks[index])
		blocks.append(iblock)

	return blocks
	
def decode(blocks):
	found = {}
	fail = 0
	prev = 0

	while blocks:

		print len(blocks)
		b = blocks.pop(0)

		if len(b.parents) == 0:
			pass

		elif len(b.parents) == 1:
			found[b.parents[0]] = b.content
		
		else:
			parents = []
			for parent in b.parents:
				if found.has_key(parent):
					b.content = xor(b.content, found[parent])
				else:
					parents.append(parent)
			b.parents = parents
			blocks.append(b)
					
	return found
	
def xor(s1, s2):
	result = ""
	for i in range(len(s1)):
		if s1[i] == s2[i]:
			result += "0"
		else:
			result += "1"
	return result
	
if __name__ == '__main__':
	f = open("sample.txt", "r")
	e = encode(f.read(), 10, 64, 10)
	d = decode(e)
	print d
	f.close()
	