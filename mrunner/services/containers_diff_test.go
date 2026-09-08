package mservices

import (
	"slices"
	"testing"

	"github.com/Liphium/magic/v4/mconfig"
	"github.com/jinzhu/copier"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/api/types/network"
)

func TestConfigMatches(t *testing.T) {
	alloc := mconfig.ContainerAllocation{Name: "mgc-app-dev-postgres", Ports: []uint{5432, 5433}}
	base := ManagedContainerOptions{
		Image: "postgres:17",
		Env:   []string{"POSTGRES_PASSWORD=x", "POSTGRES_USER=y"},
		Ports: []string{"5432/tcp", "5433/tcp"},
		Volumes: []ContainerVolume{
			{NameSuffix: "data", Target: "/var/lib/postgresql"},
		},
	}

	exposed := mustExposedSet(t, alloc, base.Ports)
	expectedBindings := mustPortBindings(t, alloc, base.Ports)

	reuse := &container.InspectResponse{
		ID: "abc",
		Config: &container.Config{
			Image:        "postgres:17",
			Env:          []string{"POSTGRES_USER=y", "POSTGRES_PASSWORD=x"}, // unsorted on purpose
			ExposedPorts: exposed,
			Cmd:          []string{"some", "start", "command"},
		},
		HostConfig: &container.HostConfig{
			PortBindings: expectedBindings,
			Mounts: []mount.Mount{
				{Type: mount.TypeVolume, Source: "mgc-app-dev-postgres-data", Target: "/var/lib/postgresql"},
			},
		},
	}

	tests := []struct {
		name   string
		mutate func(o *ManagedContainerOptions, e *container.InspectResponse)
		want   bool
	}{
		{name: "identical", want: true},
		{
			name:   "image changed",
			mutate: func(o *ManagedContainerOptions, _ *container.InspectResponse) { o.Image = "postgres:18" },
			want:   false,
		},
		{
			name: "command changed",
			mutate: func(o *ManagedContainerOptions, e *container.InspectResponse) {
				o.Cmd = []string{
					"some", "other", "command",
				}
			},
			want: false,
		},
		{
			name: "command changed by container image",
			mutate: func(o *ManagedContainerOptions, e *container.InspectResponse) {
				o.Cmd = nil
				e.Config.Cmd = []string{"some", "other", "command"}
			},
			want: true, // The driver should ignore this as no command specified
		},
		{
			name: "env changed",
			mutate: func(o *ManagedContainerOptions, _ *container.InspectResponse) {
				o.Env = append(o.Env, "POSTGRES_DB=postgres")
			},
			want: false,
		},
		{
			name: "env reordered",
			mutate: func(_ *ManagedContainerOptions, e *container.InspectResponse) {
				e.Config.Env = []string{"POSTGRES_USER=y", "POSTGRES_PASSWORD=x"}
			},
			want: true,
		},
		{
			name:   "cmd added",
			mutate: func(o *ManagedContainerOptions, _ *container.InspectResponse) { o.Cmd = []string{"server", "-s3"} },
			want:   false,
		},
		{
			name: "one port less",
			mutate: func(_ *ManagedContainerOptions, e *container.InspectResponse) {
				// Requested opts stays on 5432; existing now only exposes/binds 5433.
				allocCopy := alloc
				allocCopy.Ports = []uint{5433}
				e.Config.ExposedPorts = mustExposedSet(t, allocCopy, []string{"5433/tcp"})
				e.HostConfig.PortBindings = mustPortBindings(t, allocCopy, []string{"5433/tcp"})
			},
			want: false,
		},
		{
			name: "host port changed",
			mutate: func(_ *ManagedContainerOptions, e *container.InspectResponse) {
				allocCopy := alloc
				allocCopy.Ports = []uint{9999, 5433} // requested opts stays on 5432; existing now binds 5432 to 9999
				e.HostConfig.PortBindings = mustPortBindings(t, allocCopy, base.Ports)
			},
			want: false,
		},
		{
			name: "ports rotated",
			mutate: func(o *ManagedContainerOptions, e *container.InspectResponse) {
				allocCopy := alloc
				allocCopy.Ports = slices.Clone(allocCopy.Ports)
				slices.Reverse(allocCopy.Ports)
				baseCopy := slices.Clone(base.Ports)
				slices.Reverse(baseCopy)
				e.Config.ExposedPorts = mustExposedSet(t, allocCopy, baseCopy)
				e.HostConfig.PortBindings = mustPortBindings(t, allocCopy, baseCopy)
			},
			want: true,
		},
		{
			name: "missing volume",
			mutate: func(_ *ManagedContainerOptions, e *container.InspectResponse) {
				e.HostConfig.Mounts = []mount.Mount{}
			},
			want: false,
		},
		{
			name: "volume wrong type",
			mutate: func(_ *ManagedContainerOptions, e *container.InspectResponse) {
				e.HostConfig.Mounts[0].Type = mount.TypeBind
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := base
			existing := &container.InspectResponse{}
			if err := copier.Copy(existing, reuse); err != nil {
				t.Fatalf("couldn't copy inspect response: %v", err)
			}
			if tt.mutate != nil {
				tt.mutate(&opts, existing)
			}
			if reason := configMatches(existing, alloc, opts); tt.want != (reason == "") {
				expected := "something"
				if tt.want { // When the thing should be true, no reason is returned
					expected = "nothing"
				}
				t.Errorf("configMatches returned reason %v, but the reason should be %v", reason, expected)
			}
		})
	}
}

func mustExposedSet(t *testing.T, a mconfig.ContainerAllocation, ports []string) network.PortSet {
	t.Helper()
	exposed, _, err := buildPortBindings(a, ports)
	if err != nil {
		t.Fatal(err)
	}
	return exposed
}

func mustPortBindings(t *testing.T, a mconfig.ContainerAllocation, ports []string) network.PortMap {
	t.Helper()
	_, bindings, err := buildPortBindings(a, ports)
	if err != nil {
		t.Fatal(err)
	}
	return bindings
}
