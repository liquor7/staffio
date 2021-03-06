package backends

import (
	"log"

	"github.com/liut/staffio/pkg/backends/ldap"
	"github.com/liut/staffio/pkg/models"
	"github.com/liut/staffio/pkg/models/cas"
	"github.com/liut/staffio/pkg/models/common"
	"github.com/liut/staffio/pkg/settings"
)

type Servicer interface {
	models.Authenticator
	models.StaffStore
	models.PasswordStore
	models.GroupStore
	cas.TicketStore
	OSIN() OSINStore
	CloseAll()
	StoreStaff(*models.Staff) error
	InGroup(gn, uid string) bool
	ProfileModify(uid, password string, staff *models.Staff) error
	PasswordForgot(at common.AliasType, target, uid string) error
	PasswordResetTokenVerify(token string) (uid string, err error)
	PasswordResetWithToken(login, token, passwd string) (err error)
}

type serviceImpl struct {
	*ldap.LDAPStore
	osinStore *DbStorage
}

var _ Servicer = (*serviceImpl)(nil)

func NewService() Servicer {

	cfg := &ldap.Config{
		Addr:   settings.LDAP.Hosts,
		Base:   settings.LDAP.Base,
		Bind:   settings.LDAP.BindDN,
		Passwd: settings.LDAP.Password,
	}
	store, err := ldap.NewStore(cfg)
	if err != nil {
		log.Fatalf("new service ERR %s", err)
	}
	// LDAP is a special store
	return &serviceImpl{
		LDAPStore: store,
		osinStore: NewStorage(),
	}

}

func (s *serviceImpl) OSIN() OSINStore {
	return s.osinStore
}

func (s *serviceImpl) CloseAll() {
	s.LDAPStore.Close()
	s.osinStore.Clone()
}
