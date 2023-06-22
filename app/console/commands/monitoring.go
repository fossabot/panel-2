package commands

import (
	"strconv"

	"github.com/gookit/color"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"

	"panel/app/models"
	"panel/packages/helpers"
)

type Monitoring struct {
}

// Signature The name and signature of the console command.
func (receiver *Monitoring) Signature() string {
	return "panel:monitoring"
}

// Description The console command description.
func (receiver *Monitoring) Description() string {
	return "[面板] 系统监控"
}

// Extend The console command extend.
func (receiver *Monitoring) Extend() command.Extend {
	return command.Extend{
		Category: "panel",
	}
}

// Handle Execute the console command.
func (receiver *Monitoring) Handle(ctx console.Context) error {
	var setting models.Setting
	if err := facades.Orm().Query().Where("key", "monitor").First(&setting); err != nil {
		return nil
	}
	if setting.Value == "0" || len(setting.Value) == 0 {
		return nil
	}

	info, err := helpers.GetMonitoringInfo()
	if err != nil {
		facades.Log().Errorf("[面板] 系统监控失败: %s", err.Error())
		color.Redf("[面板] 系统监控失败: %s", err.Error())
		return nil
	}

	err = facades.Orm().Query().Create(&models.Monitor{
		Info: info,
	})
	if err != nil {
		facades.Log().Errorf("[面板] 系统监控保存失败: %s", err.Error())
		color.Redf("[面板] 系统监控保存失败: %s", err.Error())
		return nil
	}

	// 删除过期数据
	err = facades.Orm().Query().Where("key", "monitor_days").First(&setting)
	if err != nil {
		return nil
	}
	if setting.Value == "0" || len(setting.Value) == 0 {
		return nil
	}

	days, err := strconv.Atoi(setting.Value)
	if err != nil {
		return nil
	}
	_, err = facades.Orm().Query().Where("created_at < ?", carbon.Now().SubDays(days).ToDateTimeString()).Delete(&models.Monitor{})
	if err != nil {
		facades.Log().Errorf("[面板] 系统监控删除过期数据失败: %s", err.Error())
		return nil
	}

	return nil
}