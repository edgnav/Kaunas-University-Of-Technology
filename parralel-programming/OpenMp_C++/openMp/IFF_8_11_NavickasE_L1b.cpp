#include <iostream>
#include <iomanip>
#include <sstream>
#include <fstream>
#include <string>
#include <nlohmann/json.hpp>
#include <omp.h>
#include <ostream>
#include <sha1.hpp>
#include "TextTable.h"

using namespace std;
using json = nlohmann::json;




struct Musician {
	string name="";
	int careerLength=0;
	float salary=0;
	string hash="";
	
};

struct DataMonitor {
	Musician Musicians[10];
	omp_lock_t* Lock;
	bool DoneOrNOt = false;
	int CurrentSize = 0;
	int To = 0;
	int From = 0;
	bool Insert(Musician musician) {
		omp_set_lock(Lock);
		if (CurrentSize < 10) {
			Musicians[To] = musician;
			To = (To + 1) % 10;
			CurrentSize++;
			omp_unset_lock(Lock);
			return true;
		}
		else {
			omp_unset_lock(Lock);
			return false;
		}
	}
	bool Remove(Musician* musician) {
		omp_set_lock(Lock);
		if (CurrentSize > 0) {
			*musician = Musicians[From];
			Musician emptyMusician = Musician();
			Musicians[From] = emptyMusician;
			From = (From + 1) % 10;
			CurrentSize--;
			omp_unset_lock(Lock);
			return true;
		}
		else {
			if (DoneOrNOt && CurrentSize == 0)
			{
				omp_unset_lock(Lock);
				return true;

			}
			omp_unset_lock(Lock);
			return false;
		}
	}
};
struct ResultsMonitor {
	Musician Musicians[35];
	omp_lock_t* Lock;
	bool IsDone = false;
	int CurrentSize = 0;
	int To = 0;
	int From = 0;
	void InsertAndSortData(Musician musician) {
		omp_set_lock(Lock);

			int i = 25;
			while (musician.hash > Musicians[i].hash) {
				i--;
				if (i < 0) {
					break;
				}
				Musicians[i + 1] = Musicians[i];
			}
			Musicians[i + 1] = musician;
			CurrentSize++;
		
		omp_unset_lock(Lock);
	}
};
string hashMusicians(Musician musician) {
	stringstream stream;
	stream << "{" << musician.name << " " << musician.careerLength << " " << fixed << setprecision(2) << musician.salary << " }";
	string s = stream.str();
	int size = s.size();
	const char* chars = s.c_str();
	char hex[41];
	string hash = "";
	sha1(chars).finalize().print_hex(hex);
	cout << hash << endl;
	for (int i = 0; i < 40; i++) {
		hash += hex[i];
	}
	

	//cout << hash <<  endl;

	return hash;
}
void WriteToFile(vector<Musician> initialData, Musician Musicians[25]) {


	TextTable t('-', '|', '+');
	t.add("NR ");
	t.add("Name          ");
	t.add("Career length ");
	t.add("Salary   ");
	t.endOfRow();
	int i = 1;
	for (Musician musician : initialData) {
		t.add(to_string(i));
		t.add(musician.name);
		t.add(to_string(musician.careerLength));
		t.add(to_string(musician.salary));
		t.endOfRow();
		i++;

	}
	t.setAlignment(2, TextTable::Alignment::RIGHT);
	cout << t << endl;

	TextTable tt('-', '|', '+');
	tt.add("NR");
	tt.add("Name");
	tt.add("Career length");
	tt.add("Salary");
	tt.add("Hash");
	tt.endOfRow();
	
	for(int i=0;i<25;i++) {
		/*if (Musicians[i].hash == "")
			continue;*/
		tt.add(to_string(i+1));
		tt.add(Musicians[i].name);
		tt.add(to_string(Musicians[i].careerLength));
		tt.add(to_string(Musicians[i].salary));
		tt.add(Musicians[i].hash);
		tt.endOfRow();

	}
	tt.setAlignment(2, TextTable::Alignment::RIGHT);
	cout << tt << endl;


	ofstream file;
	file.open("IFF-8-11_NavickasE_L1_rez.txt");
	
	file << "Initial data\n";
	file << t;
	file << "\n";
	file << "Results\n";
	file << "\n";
	file << tt;
	file.close();
}


vector<Musician> JsonToVector(json json) {
	vector<Musician> Musicians(json.size());
	for (int i = 0; i < (int)json.size(); i++) {
		Musician musician = Musician();
		musician.name = json[i].at("name");
		musician.careerLength = json[i].at("careerLength");
		musician.salary = json[i].at("salary");
		Musicians[i] = musician;
	}
	return Musicians;
}
int main() {
	
	ifstream i("IFF-8-11_EdgarasNavickas_L1_dat_3.json");
	Musician musiciansArray[25];
	
	DataMonitor x;
	
	json j;
	i >> j;
	vector<Musician> initialMusicianData = JsonToVector(j);
	omp_lock_t * lockData = new omp_lock_t;
	omp_lock_t * lockResults = new omp_lock_t;
	omp_init_lock(lockData);
	omp_init_lock(lockResults);
	DataMonitor dataMonitor = DataMonitor();
	ResultsMonitor resultsMonitor = ResultsMonitor();
	dataMonitor.Lock = lockData;
	resultsMonitor.Lock = lockResults;
#pragma omp parallel num_threads(5)
	{

		int thread_id = omp_get_thread_num();
		if (thread_id == 0) {
			for (int i = 0; i < initialMusicianData.size(); i++) {
				while (!dataMonitor.Insert(initialMusicianData[i]));
			}
			dataMonitor.DoneOrNOt = true;
		}
		else {
			while (!dataMonitor.DoneOrNOt || dataMonitor.CurrentSize != 0)
			{
				Musician musician = Musician();
				while (!dataMonitor.Remove(&musician));
					
				
				string hash = hashMusicians(musician);
				
				musician.hash = hash;
				//cout << hash << endl;
				if (isalpha(hash[0]) && isalpha(hash[39]))
				{
					
					resultsMonitor.InsertAndSortData(musician);
				}
			}
		}
	}

	omp_destroy_lock(dataMonitor.Lock);
	omp_destroy_lock(resultsMonitor.Lock);
	WriteToFile(initialMusicianData, resultsMonitor.Musicians);
	
}




