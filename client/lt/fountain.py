# Luby Transform (LT) Python Implementation for Hermes Project
# http://en.wikipedia.org/wiki/Luby_transform_code

from random import randint
from collections import Counter

class block:
	
	def __init__(self, p, c):
		self.parents = p
		self.content = c

def encode(src, dist, size, perc): # --> []blocks

	# Converts src into binary
	src = ''.join(['%08d'%int(bin(ord(i))[2:]) for i in src])
	
	# Take source and turn into chunks of size length
	# - fix trunkate if src not divisable by size
	# - dist needs to be 1 <= d <= n
	chunks = []
	for i in range(len(src) / size):
		chunks.append(src[i * size:(i * size) + size])

	# Takes chunks and combines them into blocks
	blocks = []
	for i in range(int(len(chunks) * perc)):
		index = randint(0, len(chunks)-1)
		iblock = block([index], chunks[index])
		for i in range(randint(0, dist)):
			index = randint(0, len(chunks)-1)
			iblock.parents.append(index)
			iblock.content = xor(iblock.content, chunks[index])
		blocks.append(iblock)

	return blocks
	
def decode(blocks): # --> string

	# Iterate over blocks until all are solved
	found = {}
	while blocks:

		print len(blocks)
		b = blocks.pop(0)

		# Broken case, should not exist
		if len(b.parents) == 0:
			pass

		# If block is solved (i.e. one parent)
		elif len(b.parents) == 1:
			found[b.parents[0]] = b.content
		
		# If block is not solved...
		else:

			# ... check if any parents are solved, or...
			parents = []
			for parent in b.parents:
				if found.has_key(parent):
					b.content = xor(b.content, found[parent])
				else:
					parents.append(parent)
			b.parents = parents

			# ... check if block is a subset of another
			for c in blocks:
				if not Counter(b.parents) - Counter(c.parents): # b is a subset of c
					c.content = xor(b.content, c.content)
					parents = c.parents
					for i in b.parents:
						parents.remove(i)
					c.parents = parents

			blocks.append(b)

	# Take dict found, and convert back to src string
	result = ""
	keys = found.keys()
	keys.sort()
	for i in keys:
		result += found[i]
	return ''.join(chr(int(result[i:i+8], 2)) for i in xrange(0, len(result), 8))
	
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
	src = f.read()
	f.close()
	e = encode(src, 5, 64, 2)
	d = decode(e)
	print d