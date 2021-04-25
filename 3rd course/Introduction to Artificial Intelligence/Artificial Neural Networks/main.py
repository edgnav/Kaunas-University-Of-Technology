import matplotlib.pyplot as plt
import numpy as np
from sklearn.linear_model import LinearRegression
from Data import Data

def plot(data, xlabel, ylabel, title):
    plt.title(title)
    plt.xlabel(xlabel)
    plt.ylabel(ylabel)
    plt.plot(data[0], data[1])
    plt.show()


if __name__ == "__main__":
    data = Data("sunspot.txt") 
    data.read_file()
    print(len(data.date))
    #plot(data,  "Metai", "Saulės dėmių aktyvumas", "Saulės dėmių aktyvumas")


    
    