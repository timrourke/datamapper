package datamapper

import (
	"strings"
	"testing"
)

type ModelStub struct {
	id string
}

func (m *ModelStub) GetID() string {
	return m.id
}

func errShouldContainStr(err error, s string, t *testing.T) {
	if !strings.Contains(err.Error(), s) {
		t.Errorf("error should contain string \"%s\", contains \"%s\"", s, err.Error())
	}
}

func failOnUnexpectedErr(err error, t *testing.T) {
	if err != nil {
		t.Errorf("unexpected error: \"%s\"", err.Error())
	}
}

func TestNewUnitOfWork(t *testing.T) {
	u := NewUnitOfWork()

	if len(u.newObjects) != 0 {
		t.Errorf("should have zero new objects registered\n%+v", u.newObjects)
	}

	if len(u.dirtyObjects) != 0 {
		t.Errorf("should have zero new objects registered\n%+v", u.dirtyObjects)
	}

	if len(u.removedObjects) != 0 {
		t.Errorf("should have zero new objects registered\n%+v", u.removedObjects)
	}
}

func TestAssertModelNotRegisteredAsPanicsOnUknownState(t *testing.T) {
	u := NewUnitOfWork()

	m := &ModelStub{id: "5"}

	defer func() {
		if r := recover(); r == nil {
			t.Error("should panic when unkown persistence state passed")
		}
	}()

	_ = u.assertModelNotRegisteredAs(m, "some-state")
}

func TestRegisterNew(t *testing.T) {
	m := &ModelStub{id: "5"}

	u := NewUnitOfWork()

	err := u.RegisterNew(m)
	failOnUnexpectedErr(err, t)

	foundModel, isRegisteredNew := u.newObjects["5"]

	if !isRegisteredNew {
		t.Error("should be registered as new")
	}

	if foundModel != m {
		t.Errorf(
			"registered model should be correct: \nexpected: %+v \nactual: %+v",
			m,
			foundModel,
		)
	}
}

func TestRegisterNewFailsWithNoId(t *testing.T) {
	m := &ModelStub{}

	u := NewUnitOfWork()

	err := u.RegisterNew(m)

	if err == nil {
		t.Error("should return an error if the model has no ID")
	}

	errShouldContainStr(err, "model has no ID", t)
}

func TestRegisterNewFailsIfModelDirty(t *testing.T) {
	m := &ModelStub{id: "5"}

	u := NewUnitOfWork()

	err := u.RegisterDirty(m)
	failOnUnexpectedErr(err, t)

	err = u.RegisterNew(m)

	if err == nil {
		t.Error("should return an error if the model is already registered as dirty")
	}

	errShouldContainStr(err, "already registered as dirty", t)
}

func TestRegisterNewFailsIfModelRemoved(t *testing.T) {
	m := &ModelStub{id: "5"}

	u := NewUnitOfWork()

	err := u.RegisterRemoved(m)
	failOnUnexpectedErr(err, t)

	err = u.RegisterNew(m)

	if err == nil {
		t.Error("should return an error if the model is already registered as removed")
	}

	errShouldContainStr(err, "already registered as removed", t)
}

func TestRegisterNewFailsIfModelAlreadyNew(t *testing.T) {
	m := &ModelStub{id: "5"}

	u := NewUnitOfWork()

	err := u.RegisterNew(m)
	failOnUnexpectedErr(err, t)

	err = u.RegisterNew(m)

	if err == nil {
		t.Error("should return an error if the model is already registered as new")
	}

	errShouldContainStr(err, "already registered as new", t)
}

func TestRegisterDirty(t *testing.T) {
	m := &ModelStub{id: "5"}

	u := NewUnitOfWork()

	err := u.RegisterDirty(m)
	failOnUnexpectedErr(err, t)

	foundModel, isRegisteredDirty := u.dirtyObjects["5"]

	if !isRegisteredDirty {
		t.Error("should be registered as dirty")
	}

	if foundModel != m {
		t.Errorf(
			"registered model should be correct: \nexpected: %+v \nactual: %+v",
			m,
			foundModel,
		)
	}
}

func TestRegisterDirtyFailsWithNoId(t *testing.T) {
	m := &ModelStub{}

	u := NewUnitOfWork()

	err := u.RegisterDirty(m)

	if err == nil {
		t.Error("should return an error if the model has no ID")
	}

	errShouldContainStr(err, "model has no ID", t)
}

func TestRegisterDirtyDoesNothingIfModelNew(t *testing.T) {
	m := &ModelStub{id: "5"}

	u := NewUnitOfWork()

	err := u.RegisterNew(m)
	failOnUnexpectedErr(err, t)

	foundModel, isRegisteredNew := u.newObjects[m.GetID()]
	if !isRegisteredNew {
		t.Error("model should be registered as new")
	}

	if foundModel != m {
		t.Errorf(
			"registered model should be correct: \nexpected: %+v \nactual: %+v",
			m,
			foundModel,
		)
	}

	err = u.RegisterDirty(m)
	failOnUnexpectedErr(err, t)

	_, isRegisteredDirty := u.dirtyObjects[m.GetID()]
	if isRegisteredDirty {
		t.Error("model should not be registered as dirty if already registered as new")
	}

	foundModel, isStillRegisteredNew := u.newObjects[m.GetID()]
	if !isStillRegisteredNew {
		t.Error("model should still be registered as new")
	}

	if foundModel != m {
		t.Errorf(
			"registered model should be correct: \nexpected: %+v \nactual: %+v",
			m,
			foundModel,
		)
	}
}

func TestRegisterDirtyFailsIfModelRemoved(t *testing.T) {
	m := &ModelStub{id: "5"}

	u := NewUnitOfWork()

	err := u.RegisterRemoved(m)
	failOnUnexpectedErr(err, t)

	err = u.RegisterDirty(m)

	if err == nil {
		t.Error("should return an error if the model is already registered as removed")
	}

	errShouldContainStr(err, "already registered as removed", t)
}

func TestRegisterRemoved(t *testing.T) {
	m := &ModelStub{id: "5"}

	u := NewUnitOfWork()

	err := u.RegisterRemoved(m)
	failOnUnexpectedErr(err, t)

	foundModel, isRegisteredRemoved := u.removedObjects["5"]

	if !isRegisteredRemoved {
		t.Error("should be registered as removed")
	}

	if foundModel != m {
		t.Errorf(
			"registered model should be correct: \nexpected: %+v \nactual: %+v",
			m,
			foundModel,
		)
	}
}

func TestRegisterRemovedFailsWithNoId(t *testing.T) {
	m := &ModelStub{}

	u := NewUnitOfWork()

	err := u.RegisterRemoved(m)

	if err == nil {
		t.Error("should return an error if the model has no ID")
	}

	errShouldContainStr(err, "model has no ID", t)
}

func TestRegisterRemovedDeletesNewObject(t *testing.T) {
	m := &ModelStub{id: "5"}

	u := NewUnitOfWork()

	err := u.RegisterNew(m)
	failOnUnexpectedErr(err, t)

	err = u.RegisterRemoved(m)

	_, isStillRegisteredNew := u.newObjects["5"]
	if isStillRegisteredNew {
		t.Error("should delete registration for a new object when registering that new object as removed")
	}

	_, isRegisteredRemoved := u.removedObjects["5"]
	if isRegisteredRemoved {
		t.Error("should not register an object that was previously new as slated for removal - new object has not been persisted so nothing to delete")
	}
}
