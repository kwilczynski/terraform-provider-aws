package rds

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	rdsClusterInitiateUpgradeTimeout = 5 * time.Minute

	dbClusterRoleAssociationCreatedTimeout = 5 * time.Minute
	dbClusterRoleAssociationDeletedTimeout = 5 * time.Minute

	// EventSubscriptionDeletedTimeout = 10 * time.Minute
	dbClusterActivityStreamRetryDelay      = 5 * time.Second
	dbClusterActivityStreamRetryMinTimeout = 3 * time.Second // Minimum timeout to retry fetch RDS Cluster Activity Stream Status
)

func waitEventSubscriptionCreated(conn *rds.RDS, id string, timeout time.Duration) (*rds.EventSubscription, error) {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{EventSubscriptionStatusCreating},
		Target:     []string{EventSubscriptionStatusActive},
		Refresh:    statusEventSubscription(conn, id),
		Timeout:    timeout,
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*rds.EventSubscription); ok {
		return output, err
	}

	return nil, err
}

func waitEventSubscriptionDeleted(conn *rds.RDS, id string, timeout time.Duration) (*rds.EventSubscription, error) {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{EventSubscriptionStatusDeleting},
		Target:     []string{},
		Refresh:    statusEventSubscription(conn, id),
		Timeout:    timeout,
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*rds.EventSubscription); ok {
		return output, err
	}

	return nil, err
}

func waitEventSubscriptionUpdated(conn *rds.RDS, id string, timeout time.Duration) (*rds.EventSubscription, error) {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{EventSubscriptionStatusModifying},
		Target:     []string{EventSubscriptionStatusActive},
		Refresh:    statusEventSubscription(conn, id),
		Timeout:    timeout,
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*rds.EventSubscription); ok {
		return output, err
	}

	return nil, err
}

// waitDBProxyEndpointAvailable waits for a DBProxyEndpoint to return Available
func waitDBProxyEndpointAvailable(conn *rds.RDS, id string, timeout time.Duration) (*rds.DBProxyEndpoint, error) { //nolint:unparam
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			rds.DBProxyEndpointStatusCreating,
			rds.DBProxyEndpointStatusModifying,
		},
		Target:  []string{rds.DBProxyEndpointStatusAvailable},
		Refresh: statusDBProxyEndpoint(conn, id),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*rds.DBProxyEndpoint); ok {
		return output, err
	}

	return nil, err
}

// waitDBProxyEndpointDeleted waits for a DBProxyEndpoint to return Deleted
func waitDBProxyEndpointDeleted(conn *rds.RDS, id string, timeout time.Duration) (*rds.DBProxyEndpoint, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{rds.DBProxyEndpointStatusDeleting},
		Target:  []string{},
		Refresh: statusDBProxyEndpoint(conn, id),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*rds.DBProxyEndpoint); ok {
		return output, err
	}

	return nil, err
}

func waitDBClusterRoleAssociationCreated(conn *rds.RDS, dbClusterID, roleARN string) (*rds.DBClusterRole, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{ClusterRoleStatusPending},
		Target:  []string{ClusterRoleStatusActive},
		Refresh: statusDBClusterRole(conn, dbClusterID, roleARN),
		Timeout: dbClusterRoleAssociationCreatedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*rds.DBClusterRole); ok {
		return output, err
	}

	return nil, err
}

func waitDBClusterRoleAssociationDeleted(conn *rds.RDS, dbClusterID, roleARN string) (*rds.DBClusterRole, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{ClusterRoleStatusActive, ClusterRoleStatusPending},
		Target:  []string{},
		Refresh: statusDBClusterRole(conn, dbClusterID, roleARN),
		Timeout: dbClusterRoleAssociationDeletedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*rds.DBClusterRole); ok {
		return output, err
	}

	return nil, err
}

func waitDBInstanceDeleted(conn *rds.RDS, id string, timeout time.Duration) (*rds.DBInstance, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			InstanceStatusAvailable,
			InstanceStatusBackingUp,
			InstanceStatusConfiguringEnhancedMonitoring,
			InstanceStatusConfiguringLogExports,
			InstanceStatusCreating,
			InstanceStatusDeleting,
			InstanceStatusIncompatibleParameters,
			InstanceStatusModifying,
			InstanceStatusStarting,
			InstanceStatusStopping,
			InstanceStatusStorageFull,
			InstanceStatusStorageOptimization,
		},
		Target:     []string{},
		Refresh:    statusDBInstance(conn, id),
		Timeout:    timeout,
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*rds.DBInstance); ok {
		return output, err
	}

	return nil, err
}

func waitDBClusterInstanceDeleted(conn *rds.RDS, id string, timeout time.Duration) (*rds.DBInstance, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			InstanceStatusConfiguringLogExports,
			InstanceStatusDeleting,
			InstanceStatusModifying,
		},
		Target:     []string{},
		Refresh:    statusDBInstance(conn, id),
		Timeout:    timeout,
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*rds.DBInstance); ok {
		return output, err
	}

	return nil, err
}

// waitActivityStreamStarted waits for Aurora Cluster Activity Stream to be started
func waitActivityStreamStarted(conn *rds.RDS, dbClusterIdentifier string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for RDS Cluster Activity Stream %s to become started...", dbClusterIdentifier)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{rds.ActivityStreamStatusStarting},
		Target:     []string{rds.ActivityStreamStatusStarted},
		Refresh:    statusDBClusterActivityStream(conn, dbClusterIdentifier),
		Timeout:    timeout,
		Delay:      dbClusterActivityStreamRetryDelay,
		MinTimeout: dbClusterActivityStreamRetryMinTimeout,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for RDS Cluster Activity Stream (%s) to be started: %v", dbClusterIdentifier, err)
	}
	return nil
}

// waitActivityStreamStarted waits for Aurora Cluster Activity Stream to be stopped
func waitActivityStreamStopped(conn *rds.RDS, dbClusterIdentifier string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for RDS Cluster Activity Stream %s to become stopped...", dbClusterIdentifier)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{rds.ActivityStreamStatusStopping},
		Target:     []string{rds.ActivityStreamStatusStopped},
		Refresh:    statusDBClusterActivityStream(conn, dbClusterIdentifier),
		Timeout:    timeout,
		Delay:      dbClusterActivityStreamRetryDelay,
		MinTimeout: dbClusterActivityStreamRetryMinTimeout,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for RDS Cluster Activity Stream (%s) to be stopped: %v", dbClusterIdentifier, err)
	}
	return nil
}
