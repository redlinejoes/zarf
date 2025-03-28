import type { PlaywrightTestConfig } from '@playwright/test';
// import { devices } from '@playwright/test';

/**
 * Read environment variables from file.
 * https://github.com/motdotla/dotenv
 */
// require('dotenv').config();

/**
 * See https://playwright.dev/docs/test-configuration.
 */
const config: PlaywrightTestConfig = {
	testDir: '../test/ui',
	/* This is 15 minutes jon */
	timeout: 15 * 60 * 1000,
	expect: {
		/**
		 * Maximum time expect() should wait for the condition to be met.
		 * For example in `await expect(locator).toHaveText();`
		 */
		timeout: 15 * 1000
	},
	/* Run tests in files in parallel */
	fullyParallel: false,
	/* Fail the build on CI if you accidentally left test.only in the source code. */
	forbidOnly: !!process.env.CI,
	/* Retry on CI only */
	retries: process.env.CI ? 2 : 0,
	/* Opt out of parallel tests on CI. */
	workers: process.env.CI ? 1 : undefined,
	/* Reporter to use. See https://playwright.dev/docs/test-reporters */
	reporter: 'html',
	/* Shared settings for all the projects below. See https://playwright.dev/docs/api/class-testoptions. */
	use: {
		/* Maximum time each action such as `click()` can take. Defaults to 0 (no limit). */
		actionTimeout: 0,
		/* Base URL to use in actions like `await page.goto('/')`. */
		baseURL: process.env.CI ? 'http://localhost:3333' : 'http://localhost:5173',

		/* Collect trace when retrying the failed test. See https://playwright.dev/docs/trace-viewer */
		trace: 'on-first-retry',

		screenshot: 'only-on-failure',

		video: 'on'
	},

	/* Configure projects for major browsers */
	// projects: [
	// 	{
	// 		name: 'chromium',
	// 		use: {
	// 			...devices['Desktop Chrome']
	// 		}
	// 	},

	// 	{
	// 		name: 'firefox',
	// 		use: {
	// 			...devices['Desktop Firefox']
	// 		}
	// 	}

	// {
	// 	name: 'webkit',
	// 	use: {
	// 		...devices['Desktop Safari']
	// 	}
	// }
	// ],

	/* Folder for test artifacts such as screenshots, videos, traces, etc. */
	// outputDir: 'test-results/',

	/* Run your local dev server before starting the tests */
	webServer: {
		command: process.env.CI ? 'cd ../.. && make test-built-ui' : 'npm run dev',
		port: 3333,
		reuseExistingServer: true,
		timeout: 120 * 1000
	}
};

export default config;
