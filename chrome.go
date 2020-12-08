package main

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func getPrice(url string) string {
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// run task list
	var res string
	var selector string
	selector = "#app > div > div.scrollable > div.content-overflow > div > div.info-area > section.total-payout-area.info-area__total-area > div.total-payout-area__united-info-area > div:nth-child(1) > p.total-payout-area__united-info-area__item__amount"
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(1*time.Second),
		chromedp.Text(selector, &res, chromedp.NodeVisible, chromedp.ByID),
	)
	if err != nil {
		log.Fatal(err)
	}

	return res
}

func getSpace(url string) string {
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// run task list
	var used, total, usedSelector, totalSelector string
	usedSelector = "#app > div > div.scrollable > div.content-overflow > div > div.info-area > section.info-area__chart-area > section:nth-child(2) > div > div.disk-stat-area__info-area > div:nth-child(1) > p"
	totalSelector = "#app > div > div.scrollable > div.content-overflow > div > div.info-area > section.info-area__chart-area > section:nth-child(2) > div > p.disk-stat-area__amount"
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(1*time.Second),
		chromedp.Text(usedSelector, &used, chromedp.NodeVisible, chromedp.ByID),
		chromedp.Text(totalSelector, &total, chromedp.NodeVisible, chromedp.ByID),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(used)
	log.Println(total)

	return used + "/" + total
}

func getStatus(url string) [2]string {
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// run task list
	var used, total, price, usedSelector, totalSelector, priceSelector string
	priceSelector = "#app > div > div.scrollable > div.content-overflow > div > div.info-area > section.total-payout-area.info-area__total-area > div.total-payout-area__united-info-area > div:nth-child(1) > p.total-payout-area__united-info-area__item__amount"
	usedSelector = "#app > div > div.scrollable > div.content-overflow > div > div.info-area > section.info-area__chart-area > section:nth-child(2) > div > div.disk-stat-area__info-area > div:nth-child(1) > p"
	totalSelector = "#app > div > div.scrollable > div.content-overflow > div > div.info-area > section.info-area__chart-area > section:nth-child(2) > div > p.disk-stat-area__amount"
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(1*time.Second),
		chromedp.Text(usedSelector, &used, chromedp.NodeVisible, chromedp.ByID),
		chromedp.Text(totalSelector, &total, chromedp.NodeVisible, chromedp.ByID),
		chromedp.Text(priceSelector, &price, chromedp.NodeVisible, chromedp.ByID),
	)
	if err != nil {
		return [2]string{"err", err.Error()}
	}

	log.Println([2]string{used + "/" + total, price})

	return [2]string{used + "/" + total, price}
}

func fullScreenshot(urlstr string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Sleep(2 * time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// get layout metrics
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			// force viewport emulation
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			// capture screenshot
			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}
