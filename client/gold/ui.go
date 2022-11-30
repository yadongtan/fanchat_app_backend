package main

import (
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

func (app *Config) makeUI() {

	// 获取当前黄金价格
	openPrice, currentPrice, priceChange := app.getPriceText()
	// 将价格信息放入到容器
	priceContent := container.NewGridWithColumns(3,
		openPrice,
		currentPrice,
		priceChange)
	app.PriceContainer = priceContent
	// 将窗口加入到容器

	// get toolbar
	toolBar := app.getToolBar(app.MainWindow)
	// get app tabs
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Prices", theme.HomeIcon(), canvas.NewText("Price content goes here", nil)),
		container.NewTabItemWithIcon("Holdings", theme.InfoIcon(), canvas.NewText("Holdings content goes here", nil)),
	)
	tabs.SetTabLocation(container.TabLocationTop)
	finalContent := container.NewVBox(priceContent, toolBar, tabs)

	app.MainWindow.SetContent(finalContent)

}
