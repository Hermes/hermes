# basic fountain code implementation
from random import randint

class block:
	
	def __init__(self, p, c):
		self.parents = p
		self.content = c
	
def bin_ascii(s):
	return ''.join(chr(int(s[i:i+8], 2)) for i in xrange(0, len(s), 8))

def encode(s, d, b):
	# s~ource, d~istribution, b~lock size
	s = ''.join(['%08d'%int(bin(ord(i))[2:]) for i in s])
	
	chunks = []
	for i in range(len(s) / b):
		chunks.append(s[i:i + b])

	blocks = []
	for i in range(len(chunks)):
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
	
	while len(found.items()) < len(blocks):
		for b in blocks:
			if len(b.parents) == 1:
				found[b.parents[0]] = b.content
		# for item in found, remove from s
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
	e = encode(f.read(), 10, 64)
	print decode(e)
	