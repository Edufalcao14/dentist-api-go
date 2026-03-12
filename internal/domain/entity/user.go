package entity

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Role string

const (
	RolePatient      = "patient"
	RoleDentist      = "dentist"
	RoleReceptionist = "receptionist"
)

type Specialization string

const (
	SpecializationGeneralDentist                   Specialization = "General Dentist"
	SpecializationImplantologist                   Specialization = "Implantologist"
	SpecializationEndodontist                      Specialization = "Endodontist"
	SpecializationOrthodontist                     Specialization = "Orthodontist"
	SpecializationPediatricDentist                 Specialization = "Pediatric Dentist"
	SpecializationPeriodontist                     Specialization = "Periodontist"
	SpecializationOralAndMaxillofacialSurgeon      Specialization = "Oral and Maxillofacial Surgeon"
	SpecializationCosmeticDentist                  Specialization = "Cosmetic Dentist"
	SpecializationOralRadiologist                  Specialization = "Oral Radiologist"
	SpecializationGeriatricDentist                 Specialization = "Geriatric Dentist"
	SpecializationOrofacialHarmonizationSpecialist Specialization = "Orofacial Harmonization Specialist"
)

type AppointmentStatus string

const (
	AppointmentStatusScheduled  AppointmentStatus = "scheduled"
	AppointmentStatusConfirmed  AppointmentStatus = "confirmed"
	AppointmentStatusInProgress AppointmentStatus = "in_progress"
	AppointmentStatusCompleted  AppointmentStatus = "completed"
	AppointmentStatusCancelled  AppointmentStatus = "cancelled"
	AppointmentStatusNoShow     AppointmentStatus = "no_show"
)

type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusPaid     PaymentStatus = "paid"
	PaymentStatusFailed   PaymentStatus = "failed"
	PaymentStatusRefunded PaymentStatus = "refunded"
)

type PaymentMethod string

const (
	PaymentMethodCard      PaymentMethod = "credit_card"
	PaymentMethodCash      PaymentMethod = "cash"
	PaymentMethodPix       PaymentMethod = "pix"
	PaymentMethodDebitCard PaymentMethod = "debit_card"
)

type Tenants struct {
	ID                 uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name               string         `json:"name" gorm:"not null"`
	Adress             string         `json:"adress" gorm:"not null"`
	Cnpj               string         `json:"cnpj"`
	Phone              string         `json:"phone" gorm:"not null"`
	Stripe_customer_id string         `json:"stripe_customer_id"`
	Plan_type          string         `json:"plan_type"`
	Plan_status        string         `json:"plan_status"`
	CreatedAt          time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt          gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type Users struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	External_id string         `json:"external_id"`
	TenantID    uint           `json:"tenant_id"`
	Tenant      *Tenants       `json:"tenant" gorm:"foreignKey:TenantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Role        Role           `json:"role" gorm:"not null"`
	First_name  string         `json:"first_name" gorm:"not null"`
	Last_name   string         `json:"last_name" gorm:"not null"`
	Phone       string         `json:"phone"`
	Email       string         `json:"email"`
	Birthday    time.Time      `json:"birthday" gorm:"not null"`
	CPF         string         `json:"cpf"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type Dentists struct {
	ID             uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID         uint           `json:"user_id"`
	User           *Users         `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CroNumber      string         `json:"cro_number" gorm:"type:varchar(100);not null"`
	Specialization Specialization `json:"specialization" gorm:"type:varchar(255);not null"`
	CreatedAt      time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type Patients struct {
	ID           uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	TenantID     uint           `json:"tenant_id"`
	Tenant       *Tenants       `json:"tenant" gorm:"foreignKey:TenantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UserID       uint           `json:"user_id"`
	User         *Users         `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CPF          string         `json:"cpf" gorm:"type:varchar(14)"`
	Birthdate    time.Time      `json:"birthdate" gorm:"type:date"`
	Allergies    datatypes.JSON `json:"allergies" gorm:"type:jsonb"`
	MedicalNotes datatypes.JSON `json:"medical_notes" gorm:"type:jsonb"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type Appointments struct {
	ID          uint              `json:"id" gorm:"primaryKey;autoIncrement"`
	TenantID    uint              `json:"tenant_id"`
	Tenant      *Tenants          `json:"tenant" gorm:"foreignKey:TenantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	PatientID   uint              `json:"patient_id"`
	Patient     *Patients         `json:"patient" gorm:"foreignKey:PatientID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	DentistID   uint              `json:"dentist_id"`
	Dentist     *Dentists         `json:"dentist" gorm:"foreignKey:DentistID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	DateTime    time.Time         `json:"date_time" gorm:"not null"`
	Status      AppointmentStatus `json:"status" gorm:"type:varchar(50);not null"`
	Notes       string            `json:"notes" gorm:"type:text"`
	TotalAmount float64           `json:"total_amount" gorm:"type:decimal(10,2)"`
	CreatedAt   time.Time         `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time         `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt    `json:"deleted_at" gorm:"index"`
}

type Payments struct {
	ID            uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	AppointmentID uint           `json:"appointment_id"`
	Appointment   *Appointments  `json:"appointment" gorm:"foreignKey:AppointmentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Amount        float64        `json:"amount" gorm:"type:decimal(10,2);not null"`
	Status        PaymentStatus  `json:"status" gorm:"type:varchar(50);not null"`
	PaymentMethod PaymentMethod  `json:"payment_method" gorm:"type:varchar(50)"`
	Installments  int            `json:"installments" gorm:"default:1"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type Services struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	TenantID    uint           `json:"tenant_id"`
	Tenant      *Tenants       `json:"tenant" gorm:"foreignKey:TenantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Name        string         `json:"name" gorm:"type:varchar(255);not null"`
	Description string         `json:"description" gorm:"type:varchar(255)"`
	Price       float64        `json:"price" gorm:"type:decimal(10,2);not null"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type AppointmentServices struct {
	ID            uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	AppointmentID uint           `json:"appointment_id"`
	Appointment   *Appointments  `json:"appointment" gorm:"foreignKey:AppointmentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ServiceID     uint           `json:"service_id"`
	Service       *Services      `json:"service" gorm:"foreignKey:ServiceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Notes         string         `json:"notes" gorm:"type:text"`
	Price         float64        `json:"unit_price" gorm:"type:decimal(10,2);"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
