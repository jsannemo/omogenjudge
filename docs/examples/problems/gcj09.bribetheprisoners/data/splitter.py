N = int(input())

for i in range(N):
    f = open("{:02d}.in".format(i + 1), "w")
    f.write("{}\n".format(input()))
    f.write("{}\n".format(input()))
