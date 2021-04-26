class Data:
    def __init__(self, file):
        self.date = []
        self.spots = []
        self.file = file

    def read_file(self):
        dataMatrix = [[], []]
        f = open(self.file, "r")
        lines = f.readlines()

        for line in lines:
            splitLine = line.split('\t')
            self.date.append(int(splitLine[0]))
            self.spots.append(int(splitLine[1]))
        f.close()