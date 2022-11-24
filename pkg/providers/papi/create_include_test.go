package papi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/tools"
	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/akamai/cli/pkg/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	section    = "test_section"
	contractID = "test_contract"

	expectListIncludes = func(client *mockpapi) {
		listIncludesReq := papi.ListIncludesRequest{
			ContractID: "test_contract",
		}

		includes := papi.ListIncludesResponse{
			Includes: papi.IncludeItems{
				Items: []papi.Include{
					{
						AccountID:         "test_account",
						AssetID:           "test_asset",
						ContractID:        "test_contract",
						GroupID:           "test_group",
						IncludeID:         "inc_123456",
						IncludeName:       "test_include",
						IncludeType:       papi.IncludeTypeMicroServices,
						LatestVersion:     2,
						StagingVersion:    tools.IntPtr(1),
						ProductionVersion: tools.IntPtr(1),
					},
					{
						AccountID:      "test_account_1",
						AssetID:        "test_asset_1",
						ContractID:     "test_contract",
						GroupID:        "test_group",
						IncludeID:      "inc_456789",
						IncludeName:    "test_include_1",
						IncludeType:    papi.IncludeTypeCommonSettings,
						LatestVersion:  1,
						StagingVersion: tools.IntPtr(1),
					},
				},
			},
		}

		client.On("ListIncludes", mock.Anything, listIncludesReq).Return(&includes, nil).Once()
	}

	expectGetIncludeVersion = func(client *mockpapi) {
		getIncludeVersionReq := papi.GetIncludeVersionRequest{
			ContractID: "test_contract",
			GroupID:    "test_group",
			IncludeID:  "inc_123456",
			Version:    2,
		}

		version := papi.GetIncludeVersionResponse{
			AccountID:   "test_account",
			AssetID:     "test_asset",
			ContractID:  "test_contract",
			GroupID:     "test_group",
			IncludeID:   "inc_123456",
			IncludeName: "test_include",
			IncludeType: papi.IncludeTypeMicroServices,
			IncludeVersions: papi.Versions{
				Items: []papi.IncludeVersion{
					{
						UpdatedByUser:    "test_user",
						UpdatedDate:      "2022-08-22T07:17:48Z",
						ProductionStatus: papi.VersionStatusInactive,
						Etag:             "1d8ed19bce0833a3fe93e62ae5d5579a38cc2dbe",
						ProductID:        "prd_Site_Defender",
						RuleFormat:       "v2020-11-02",
						IncludeVersion:   2,
						StagingStatus:    papi.VersionStatusInactive,
					},
				},
			},
			IncludeVersion: papi.IncludeVersion{
				UpdatedByUser:    "test_user",
				UpdatedDate:      "2022-08-22T07:17:48Z",
				ProductionStatus: papi.VersionStatusInactive,
				Etag:             "1d8ed19bce0833a3fe93e62ae5d5579a38cc2dbe",
				ProductID:        "prd_Site_Defender",
				RuleFormat:       "v2020-11-02",
				IncludeVersion:   2,
				StagingStatus:    papi.VersionStatusInactive,
			},
		}

		client.On("GetIncludeVersion", mock.Anything, getIncludeVersionReq).Return(&version, nil).Once()
	}

	getIncludeRuleTreeReq = papi.GetIncludeRuleTreeRequest{
		ContractID:     "test_contract",
		GroupID:        "test_group",
		IncludeID:      "inc_123456",
		IncludeVersion: 2,
		RuleFormat:     "v2020-11-02",
	}

	expectListIncludeActivations = func(client *mockpapi) {
		listIncludeActivationsReq := papi.ListIncludeActivationsRequest{
			ContractID: "test_contract",
			GroupID:    "test_group",
			IncludeID:  "inc_123456",
		}

		activations := papi.ListIncludeActivationsResponse{
			AccountID:  "test_account",
			ContractID: "test_contract",
			GroupID:    "test_group",
			Activations: papi.IncludeActivationsRes{
				Items: []papi.IncludeActivation{
					{
						ActivationID:       "atv_12344",
						Network:            papi.ActivationNetworkStaging,
						ActivationType:     papi.ActivationTypeActivate,
						Status:             papi.ActivationStatusActive,
						SubmitDate:         "2022-10-27T12:27:54Z",
						UpdateDate:         "2022-10-27T12:28:54Z",
						Note:               "test staging activation",
						NotifyEmails:       []string{"test@example.com"},
						FMAActivationState: "steady",
						FallbackInfo: &papi.ActivationFallbackInfo{
							FastFallbackAttempted:      false,
							FallbackVersion:            1,
							CanFastFallback:            false,
							SteadyStateTime:            1666873734,
							FastFallbackExpirationTime: 1666877334,
						},
						IncludeID:      "inc_123456",
						IncludeName:    "test_include",
						IncludeType:    papi.IncludeTypeMicroServices,
						IncludeVersion: 1,
					},
					{
						ActivationID:       "atv_12343",
						Network:            papi.ActivationNetworkProduction,
						ActivationType:     papi.ActivationTypeActivate,
						Status:             papi.ActivationStatusActive,
						SubmitDate:         "2022-10-27T12:27:54Z",
						UpdateDate:         "2022-10-27T12:28:54Z",
						Note:               "test production activation",
						NotifyEmails:       []string{"test@example.com", "test1@example.com"},
						FMAActivationState: "steady",
						FallbackInfo: &papi.ActivationFallbackInfo{
							FastFallbackAttempted:      false,
							FallbackVersion:            1,
							CanFastFallback:            false,
							SteadyStateTime:            1666873734,
							FastFallbackExpirationTime: 1666877334,
						},
						IncludeID:      "inc_123456",
						IncludeName:    "test_include",
						IncludeType:    papi.IncludeTypeMicroServices,
						IncludeVersion: 1,
					},
					{
						ActivationID:       "atv_12342",
						Network:            papi.ActivationNetworkProduction,
						ActivationType:     papi.ActivationTypeDeactivate,
						Status:             papi.ActivationStatusActive,
						SubmitDate:         "2022-09-27T12:27:54Z",
						UpdateDate:         "2022-09-27T12:28:54Z",
						Note:               "test production deactivation",
						NotifyEmails:       []string{"test@example.com", "test1@example.com"},
						FMAActivationState: "steady",
						FallbackInfo: &papi.ActivationFallbackInfo{
							FastFallbackAttempted:      false,
							FallbackVersion:            1,
							CanFastFallback:            false,
							SteadyStateTime:            1666873734,
							FastFallbackExpirationTime: 1666877334,
						},
						IncludeID:      "inc_123456",
						IncludeName:    "test_include",
						IncludeType:    papi.IncludeTypeMicroServices,
						IncludeVersion: 1,
					},
					{
						ActivationID:       "atv_12341",
						Network:            papi.ActivationNetworkProduction,
						ActivationType:     papi.ActivationTypeActivate,
						Status:             papi.ActivationStatusActive,
						SubmitDate:         "2022-08-27T12:27:54Z",
						UpdateDate:         "2022-08-27T12:28:54Z",
						Note:               "test production old activation",
						NotifyEmails:       []string{"test@example.com", "test1@example.com"},
						FMAActivationState: "steady",
						FallbackInfo: &papi.ActivationFallbackInfo{
							FastFallbackAttempted:      false,
							FallbackVersion:            1,
							CanFastFallback:            false,
							SteadyStateTime:            1666873734,
							FastFallbackExpirationTime: 1666877334,
						},
						IncludeID:      "inc_123456",
						IncludeName:    "test_include",
						IncludeType:    papi.IncludeTypeMicroServices,
						IncludeVersion: 1,
					},
				},
			},
		}

		client.On("ListIncludeActivations", mock.Anything, listIncludeActivationsReq).Return(&activations, nil).Once()
	}

	expectAllProcessTemplates = func(p *mockProcessor, testData TFData) *mock.Call {

		call := p.On(
			"ProcessTemplates",
			testData,
		)
		return call.Return(nil)
	}
)

func TestCreateInclude(t *testing.T) {
	tests := map[string]struct {
		init                func(*mockpapi, *mockProcessor, string)
		includeName         string
		dir                 string
		snippetFilesToCheck []string
		jsonDir             string
		withError           error
	}{
		"include basic": {
			init: func(c *mockpapi, p *mockProcessor, dir string) {
				expectListIncludes(c)
				expectGetIncludeVersion(c)

				// Rule Tree
				var ruleResponse papi.GetIncludeRuleTreeResponse
				rules, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", dir, "mock_rules.json"))
				assert.NoError(t, err)
				assert.NoError(t, json.Unmarshal(rules, &ruleResponse))
				c.On("GetIncludeRuleTree", mock.Anything, getIncludeRuleTreeReq).Return(&ruleResponse, nil).Once()

				expectListIncludeActivations(c)
				expectAllProcessTemplates(p, getTestData("include basic"))
			},

			includeName: "test_include",
			dir:         "include_basic",
			jsonDir:     "include_basic/property-snippets",
			snippetFilesToCheck: []string{
				"test_include.json",
				"Content_Compression.json",
				"Static_Content.json",
				"Dynamic_Content.json",
			},
		},
		"error include not found": {
			init: func(c *mockpapi, p *mockProcessor, dir string) {
				c.On("ListIncludes", mock.Anything, papi.ListIncludesRequest{ContractID: "test_contract"}).
					Return(nil, fmt.Errorf("oops")).Once()
			},
			withError: ErrIncludeNotFound,
		},
		"error fetching include version": {
			init: func(c *mockpapi, p *mockProcessor, dir string) {
				expectListIncludes(c)
				c.On("GetIncludeVersion", mock.Anything, papi.GetIncludeVersionRequest{
					ContractID: "test_contract",
					GroupID:    "test_group",
					IncludeID:  "inc_123456",
					Version:    2,
				}).Return(nil, fmt.Errorf("oops")).Once()
			},
			withError:   ErrFetchingLatestIncludeVersion,
			includeName: "test_include",
		},
		"error include rules not found": {
			init: func(c *mockpapi, p *mockProcessor, dir string) {
				expectListIncludes(c)
				expectGetIncludeVersion(c)
				c.On("GetIncludeRuleTree", mock.Anything, getIncludeRuleTreeReq).Return(nil, fmt.Errorf("oops")).Once()
			},
			withError:   ErrIncludeRulesNotFound,
			includeName: "test_include",
		},
		"error fetching activations": {
			init: func(c *mockpapi, p *mockProcessor, dir string) {
				expectListIncludes(c)
				expectGetIncludeVersion(c)

				// Rule Tree
				var ruleResponse papi.GetIncludeRuleTreeResponse
				rules, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", dir, "mock_rules.json"))
				assert.NoError(t, err)
				assert.NoError(t, json.Unmarshal(rules, &ruleResponse))
				c.On("GetIncludeRuleTree", mock.Anything, getIncludeRuleTreeReq).Return(&ruleResponse, nil).Once()

				c.On("ListIncludeActivations", mock.Anything, papi.ListIncludeActivationsRequest{
					ContractID: "test_contract",
					GroupID:    "test_group",
					IncludeID:  "inc_123456",
				}).Return(nil, fmt.Errorf("oops")).Once()
			},
			withError:   ErrFetchingActivations,
			includeName: "test_include",
			dir:         "include_basic",
		},
		"error saving files": {
			init: func(c *mockpapi, p *mockProcessor, dir string) {
				expectListIncludes(c)
				expectGetIncludeVersion(c)

				// Rule Tree
				var ruleResponse papi.GetIncludeRuleTreeResponse
				rules, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", dir, "mock_rules.json"))
				assert.NoError(t, err)
				assert.NoError(t, json.Unmarshal(rules, &ruleResponse))
				c.On("GetIncludeRuleTree", mock.Anything, getIncludeRuleTreeReq).Return(&ruleResponse, nil).Once()

				expectListIncludeActivations(c)

				p.On("ProcessTemplates", getTestData("include basic")).Return(fmt.Errorf("oops")).Once()
			},
			withError:   ErrSavingFiles,
			includeName: "test_include",
			dir:         "include_basic",
			jsonDir:     "include_basic/property-snippets",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mc := new(mockpapi)
			mp := new(mockProcessor)
			test.init(mc, mp, test.dir)
			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))
			err := createInclude(ctx, contractID, test.includeName, section, fmt.Sprintf("./testdata/res/%s", test.jsonDir), "./", mc, mp)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "expected: %s; got: %s", test.withError, err)
				return
			}
			for _, f := range test.snippetFilesToCheck {
				expected, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", test.jsonDir, f))
				require.NoError(t, err)
				result, err := os.ReadFile(fmt.Sprintf("./testdata/res/%s/%s", test.jsonDir, f))
				require.NoError(t, err)
				assert.Equal(t, string(expected), string(result))
			}
			require.NoError(t, err)
			mc.AssertExpectations(t)
			mp.AssertExpectations(t)
		})
	}
}

func TestProcessIncludeTemplates(t *testing.T) {
	tests := map[string]struct {
		givenData    TFData
		dir          string
		filesToCheck []string
	}{
		"include basic": {
			givenData:    getTestData("include basic"),
			dir:          "include_basic",
			filesToCheck: []string{"includes.tf", "variables.tf", "import.sh"},
		},
		"include single network": {
			givenData:    getTestData("include single network"),
			dir:          "include_single_network",
			filesToCheck: []string{"includes.tf", "variables.tf", "import.sh"},
		},
		"include no network": {
			givenData:    getTestData("include no network"),
			dir:          "include_no_network",
			filesToCheck: []string{"includes.tf", "variables.tf", "import.sh"},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s", test.dir), 0755))
			processor := templates.FSTemplateProcessor{
				TemplatesFS: templateFiles,
				TemplateTargets: map[string]string{
					"includes.tmpl":  fmt.Sprintf("./testdata/res/%s/includes.tf", test.dir),
					"variables.tmpl": fmt.Sprintf("./testdata/res/%s/variables.tf", test.dir),
					"imports.tmpl":   fmt.Sprintf("./testdata/res/%s/import.sh", test.dir),
				},
			}
			require.NoError(t, processor.ProcessTemplates(test.givenData))

			for _, f := range test.filesToCheck {
				expected, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", test.dir, f))
				require.NoError(t, err)
				result, err := os.ReadFile(fmt.Sprintf("./testdata/res/%s/%s", test.dir, f))
				require.NoError(t, err)
				assert.Equal(t, string(expected), string(result))
			}
		})
	}

}

func getTestData(key string) TFData {
	TFDataMap := map[string]TFData{
		"include basic": {
			Section: section,
			Includes: []TFIncludeData{
				{
					ActivationNoteProduction:   "test production activation",
					ActivationNoteStaging:      "test staging activation",
					ContractID:                 "test_contract",
					ActivationEmailsProduction: []string{"test@example.com", "test1@example.com"},
					ActivationEmailsStaging:    []string{"test@example.com"},
					GroupID:                    "test_group",
					IncludeID:                  "inc_123456",
					IncludeName:                "test_include",
					IncludeType:                string(papi.IncludeTypeMicroServices),
					Networks:                   []string{"staging", "production"},
					ProductID:                  "prd_Site_Defender",
					RuleFormat:                 "v2020-11-02",
					VersionProduction:          "1",
					VersionStaging:             "1",
				},
			},
		},
		"include single network": {
			Section: section,
			Includes: []TFIncludeData{
				{
					ContractID:              "test_contract",
					ActivationEmailsStaging: []string{"test@example.com"},
					GroupID:                 "test_group",
					IncludeID:               "inc_123456",
					IncludeName:             "test_include",
					IncludeType:             string(papi.IncludeTypeMicroServices),
					Networks:                []string{"staging"},
					ProductID:               "prd_Site_Defender",
					RuleFormat:              "v2020-11-02",
					VersionStaging:          "3",
				},
			},
		},
		"include no network": {
			Section: section,
			Includes: []TFIncludeData{
				{
					ContractID:  "test_contract",
					GroupID:     "test_group",
					IncludeID:   "inc_123456",
					IncludeName: "test_include",
					IncludeType: string(papi.IncludeTypeMicroServices),
					ProductID:   "prd_Site_Defender",
					RuleFormat:  "v2020-11-02",
				},
			},
		},
	}

	return TFDataMap[key]
}
