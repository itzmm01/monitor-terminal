import sys

if len(sys.argv) == 2:
    print(eval(sys.argv[1]))
else:
    print(sys.argv[0], '"1 > 2"')
    sys.exit(1)