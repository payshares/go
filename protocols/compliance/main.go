package compliance

// AuthStatus represents auth status returned by Auth Server
type AuthStatus string

const (
	// AuthStatusOk is returned when authentication was successful
	AuthStatusOk AuthStatus = "ok"
	// AuthStatusPending is returned when authentication is pending
	AuthStatusPending AuthStatus = "pending"
	// AuthStatusDenied is returned when authentication was denied
	AuthStatusDenied AuthStatus = "denied"
)

// AuthRequest represents auth request sent to compliance server
type AuthRequest struct {
	// Marshalled AuthData JSON object (because of the attached signature, json can be marshalled to infinite number of valid JSON strings)
	DataJSON string `name:"data" valid:"required,json"`
	// Signature of sending FI
	Signature string `name:"sig" valid:"required,base64"`
}

// AuthData represents how AuthRequest.Data field looks like.
type AuthData struct {
	// The payshares address of the customer that is initiating the send.
	Sender string `json:"sender" valid:"required,payshares_address"`
	// If the caller needs the recipient's AML info in order to send the payment.
	NeedInfo bool `json:"need_info" valid:"-"`
	// The transaction that the sender would like to send in XDR format. This transaction is unsigned.
	Tx string `json:"tx" valid:"required,base64"`
	// The full text of the attachment the hash of this attachment is included in the transaction.
	AttachmentJSON string `json:"attachment" valid:"required,json"`
}

// AuthResponse represents response sent by auth server
type AuthResponse struct {
	// If this FI is willing to share AML information or not. {ok, denied, pending}
	InfoStatus AuthStatus `json:"info_status"`
	// If this FI is willing to accept this transaction. {ok, denied, pending}
	TxStatus AuthStatus `json:"tx_status"`
	// (only present if info_status is ok) JSON of the recipient's AML information. in the Payshares attachment convention
	DestInfo string `json:"dest_info,omitempty"`
	// (only present if info_status or tx_status is pending) Estimated number of seconds till the sender can check back for a change in status. The sender should just resubmit this request after the given number of seconds.
	Pending int `json:"pending,omitempty"`
}

// Attachment represents preimage object of compliance protocol in
// Payshares attachment convention
type Attachment struct {
	Nonce       string `json:"nonce"`
	Transaction `json:"transaction"`
	Operations  []Operation `json:"operations"`
}

// Transaction represents transaction field in Payshares attachment
type Transaction struct {
	SenderInfo map[string]string `json:"sender_info"`
	Route      Route             `json:"route"`
	Note       string            `json:"note"`
	Extra      string            `json:"extra"`
}

// Operation represents a single operation object in Payshares attachment
type Operation struct {
	// Overriddes Transaction field for this operation
	SenderInfo map[string]string `json:"sender_info"`
	// Overriddes Transaction field for this operation
	Route Route `json:"route"`
	// Overriddes Transaction field for this operation
	Note string `json:"note"`
}

// Route allows unmarshalling both integer and string types into string
type Route string

// SenderInfo is a helper structure with standardized fields that contains
// information about the sender. Use Map() method to transform it to
// map[string]string used in Transaction and Operation structs.
type SenderInfo struct {
	FirstName   string `json:"first_name,omitempty"`
	MiddleName  string `json:"middle_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Address     string `json:"address,omitempty"`
	City        string `json:"city,omitempty"`
	Province    string `json:"province,omitempty"`
	PostalCode  string `json:"postal_code,omitempty"`
	Country     string `json:"country,omitempty"`
	Email       string `json:"email,omitempty"`
	Phone       string `json:"phone,omitempty"`
	DateOfBirth string `json:"date_of_birth,omitempty"`
	CompanyName string `json:"company_name,omitempty"`
}
