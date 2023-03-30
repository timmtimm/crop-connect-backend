package treatment_records

import (
	"marketplace-backend/dto"
	"mime/multipart"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Domain struct {
	ID           primitive.ObjectID
	RequesterID  primitive.ObjectID
	AccepterID   primitive.ObjectID
	BatchID      primitive.ObjectID
	Number       int
	Date         primitive.DateTime
	Status       string
	Description  string
	Treatment    []dto.Treatment
	RevisionNote string
	WarningNote  string
	CreatedAt    primitive.DateTime
	UpdatedAt    primitive.DateTime
}

type Query struct {
	Skip      int64
	Limit     int64
	Sort      string
	Order     int
	FarmerID  primitive.ObjectID
	Commodity string
	Batch     string
	Number    int
	Status    string
}

// timestone
// validate dari sisi validator => kalau ada revisi, itu perlu diselesain dulu sama si farmer
// coba dicek lagi kalau revisi itu kondisinya dibuat waiting reponse aja atau engga (not solved) => solusi dibuat revisi aja biar
// si validator tau kalo ini nunggu revisi dari si petani

// kalau udah accept baru bisa nambah lagi (ini udah keimplement) dari validator (buat ngajuin lagi ya)
// kalau udah accept si petani baru bisa ngajuin harvestnya

// ini nanti nambah lagi semisal udah harvest, si validator gak bisa ngajuin catatan perawtannya lagi
// ini coba dipikiran lagi kalau ternyata si validator itu ngajuin pengisian tetapi ternyata sudah panen => sistemnya gimana apakah si validator
// bisa ngapus si permintaannya atau engga (not solved)

type Repository interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetNewestByBatchID(batchID primitive.ObjectID) (Domain, error)
	CountByBatchID(batchID primitive.ObjectID) (int, error)
	GetByID(id primitive.ObjectID) (Domain, error)
	GetByBatchID(batchID primitive.ObjectID) ([]Domain, error)
	GetByQuery(query Query) ([]Domain, int, error)
	// Update
	Update(domain *Domain) (Domain, error)
	// Delete
}

type UseCase interface {
	// Create
	RequestToFarmer(domain *Domain) (Domain, int, error)
	// Read
	GetByPaginationAndQuery(query Query) ([]Domain, int, int, error)
	GetByBatchID(batchID primitive.ObjectID) ([]Domain, int, error)
	// Update
	FillTreatmentRecord(domain *Domain, farmerID primitive.ObjectID, images []*multipart.FileHeader, notes []string) (Domain, int, error)
	Validate(domain *Domain, validatorID primitive.ObjectID) (Domain, int, error)
	UpdateNotes(domain *Domain) (Domain, int, error)
	// Delete
}
