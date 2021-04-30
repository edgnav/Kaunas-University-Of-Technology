import matplotlib.pyplot as plt
import numpy as np
from sklearn.linear_model import LinearRegression
from Data import Data

def plot(data, xlabel, ylabel, title):
    plt.title(title)
    plt.xlabel(xlabel)
    plt.ylabel(ylabel)
    plt.plot(data.date, data.spots)
    plt.show()

def plotError(data, xlabel, ylabel, title):
    plt.title(title)
    plt.xlabel(xlabel)
    plt.ylabel(ylabel)
    plt.plot(data[0], data[1])
    plt.show()

def plotErrorHist(title, xlabel, ylabel, testError):
    plt.title(title)
    plt.xlabel(xlabel)
    plt.ylabel(ylabel)
    plt.hist(testError)
    plt.show()

def split_data(sunspots, n):
    P = []
    T = []

    for i in range(len(sunspots)-n):
        temp = []
        for j in range(i, i + n):
            temp.append(sunspots[j])
        P.append(temp)
        T.append(sunspots[i+n])
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

def plotComparison(years, Tu, Tsu, xlabel, ylabel, title):
    plt.title(title)
    plt.xlabel(xlabel)
    plt.ylabel(ylabel)
    plt.plot(years, Tsu, 'red', label="Predictions")
    plt.plot(years, Tu, 'blue', label="True values")
    plt.legend()
    plt.show()


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
    #7
    Lu = 200 #amount of training data 
    PuTrain = np.array(P[:Lu])
    TuTrain = np.array(T[:Lu])
    #8
    autoregressionModel = LinearRegression().fit(PuTrain, TuTrain)
    #9
    print(autoregressionModel.coef_)
    #r_sq = autoregressionModel.score(PuTrain, TuTrain)
    #print('coefficient of determination:', r_sq)

    #10
    #train data
    TsuTrain = autoregressionModel.predict(PuTrain)

    #plotComparison(data.date[0:Lu],TuTrain,TsuTrain,"Metai","Saulės dėmių skaičius","Testavimo ir prognozavimo duomenų palyginimas su apmokymo duomenimis")

    #test data
    PuTest = np.array(P[Lu:])
    TuTest = np.array(T[Lu:])
    TsuTest = autoregressionModel.predict(PuTest)

    
    length = len(PuTest)+Lu #PuTest array is shorter, thats why we need to describe new length
    #plotComparison(data.date[Lu:length],TuTest,TsuTest,"Metai","Saulės dėmių skaičius","Testavimo ir prognozavimo duomenų palyginimas su nematytais duomenimis")

    #11
    testError = TuTest - TsuTest
    
    #plotError([data.date[200:length], testError], "Prognozės klaida", "Metai", "Klaidos dydis")

    #12
    plotErrorHist("Klaidų dydžių histograma", "Klaidos dydis", "Dažnis", testError)
    
    #13

    sum=0
    for i in range(1, len(testError)):
        sum += (testError[i])**2
    MSE = (1/length)*sum
    print("MSE = ",MSE)

    MAD=np.median(abs(testError))
    print("MAD = ",MAD)

    #14
    


    


    
    