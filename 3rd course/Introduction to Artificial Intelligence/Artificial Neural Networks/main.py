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

def split_data(sunspots, n):
    P = []
    T = []
    a = [1,2,3,4,5,6,7,8,9]

    for i in range(len(sunspots)-n):
        temp = []
        for j in range(i, i + n):
            temp.append(sunspots[j])
        P.append(temp)
        T.append(i+n)
    return P, T

def plot_3D(P, T):
    fig = plt.figure()
    ax = fig.add_subplot(projection='3d')
    ax.set_title("Įvesties ir išvesties sąrašų atvaizdavimas")
    ax.set_xlabel('Pirmos įvesties saulės dėmių skaičius')
    ax.set_ylabel('Antros įvesties saulės dėmių skaičius')
    ax.set_zlabel('Išvesties saulės dėmių skaičius')
    axis1 = []
    axis2 = []
    for i in range(len(P)):
        axis1.append(P[i][0])
        axis2.append(P[i][1])

    ax.scatter(axis1,axis2,T)
    plt.show();    


if __name__ == "__main__":
    data = Data("sunspot.txt") 
    #2--------------------------------
    data.read_file()
    #4--------------------------------
    #plot(data,  "Metai", "Saulės dėmių aktyvumas", "Saulės dėmių aktyvumas")
    #5--------------------------------
    P, T = split_data(data.spots,2)
    #6--------------------------------
    #plot_3D(P,T)
    #model = LinearRegression() #creating an instance of the class LinearRegression

    


    
    