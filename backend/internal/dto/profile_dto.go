package dto

type CompleteProfileRequest struct {
	NickName             string   `json:"nickName"`
	Department           string   `json:"department"`
	Location             string   `json:"location"`
	Designation          string   `json:"designation"`
	JobRole              string   `json:"jobRole"`
	EmploymentType       string   `json:"employmentType"`
	EmployeeStatus       string   `json:"employeeStatus"`
	SourceOfHire         string   `json:"sourceOfHire"`
	DateOfJoining        string   `json:"dateOfJoining"`
	DateOfBirth          string   `json:"dateOfBirth"`
	MaritalStatus        string   `json:"maritalStatus"`
	AboutMe              string   `json:"aboutMe"`
	Expertise            string   `json:"expertise"`
	UAN                  string   `json:"uan"`
	PAN                  string   `json:"pan"`
	WorkPhoneNumber      string   `json:"workPhoneNumber"`
	PersonalMobileNumber string   `json:"personalMobileNumber"`
	Extension            string   `json:"extension"`
	PersonalEmailAddress string   `json:"personalEmailAddress"`
	SeatingLocation      string   `json:"seatingLocation"`
	Tags                 []string `json:"tags"`
	PresentAddress       string   `json:"presentAddress"`
	PermanentAddress     string   `json:"permanentAddress"`
}