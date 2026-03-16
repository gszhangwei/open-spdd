---
name: /spdd-api-test
id: spdd-api-test
category: Testing
description: Generate a self-contained shell script with cURL commands to test API endpoints based on generated code and acceptance criteria
---

Generate a comprehensive, self-contained shell script (`scripts/test-api.sh`) with cURL commands to test all API scenarios defined in the codebase or acceptance criteria.

**Input**: The argument after `/spdd-api-test` is a reference to generated code, acceptance criteria document, or API specification.

Input can be provided in several ways:

1. **File/folder reference**: Using `@` to reference files containing API implementations or ACs
2. **Text description**: Direct text describing the API endpoints to test
3. **Combined**: Both file references and additional context

**Examples**:

```
# Reference to generated code
/spdd-api-test @src/main/java/com/example/controller/BillingController.java

# Reference to acceptance criteria
/spdd-api-test @spdd/prompt/GGQPA-XXX-202603131530-[Feat]-token-usage-billing.md

# Reference to multiple files
/spdd-api-test @src/controllers/ @requirements/api-spec.md

# With additional context
/spdd-api-test @src/api/ test the billing endpoints with edge cases for zero tokens
```

**Steps**

1. **Validate and consolidate input context**

   a. **If no input provided**, use the **AskUserQuestion tool** (open-ended, no preset options) to ask:
   > "Please provide the API implementation files, acceptance criteria, or API specification to generate tests for (you can use @file references or text description)."

   **IMPORTANT**: Do NOT proceed without input context.

   b. **If input contains `@` file/folder references**:
   - Read ALL referenced files completely using the Read tool
   - For folder references, read all relevant API-related files (controllers, routes, handlers)
   - Consolidate all file contents into a unified API context

   c. **Extract API information**:
   - Identify all API endpoints (HTTP method, path, request body schema)
   - Extract acceptance criteria and expected behaviors
   - Identify validation rules and error scenarios
   - Note any seed data or test data requirements

2. **Identify seed data and test fixtures**

   a. **Search for existing seed data**:
   - Look for seed SQL files, fixtures, or test data in common locations:
      - `src/main/resources/db/migration/`
      - `src/test/resources/`
      - `fixtures/`
      - `seed/`
      - `scripts/seed*.sql`
   - Read relevant seed files to understand available test data

   b. **Document available test entities**:
   - Customer IDs, names, and states
   - Plan IDs, pricing configurations
   - Quota configurations and limits
   - Any other domain-specific test data

   c. **If no seed data found**, derive test data from:
   - Entity definitions in the codebase
   - Example values in acceptance criteria
   - Common test data patterns

3. **Generate test script structure**

   Create the shell script with the following structure:

   ```bash
   #!/bin/bash
   # =============================================================================
   # API Test Script
   # Generated for: [API/Feature Name]
   # =============================================================================
   #
   # Usage: ./scripts/test-api.sh [BASE_URL]
   #        Default BASE_URL: http://localhost:8080
   #
   # Requirements:
   # - No external dependencies (no jq, only curl and bash)
   # - Each request has -m 10 timeout to prevent hanging
   # - HTTP status captured via: -o /tmp/response.txt -w "%{http_code}"
   #
   # =============================================================================

   # -----------------------------------------------------------------------------
   # CONFIGURATION
   # -----------------------------------------------------------------------------
   BASE_URL="${1:-http://localhost:8080}"

   # Colors for output (disabled if not a terminal)
   if [ -t 1 ]; then
       RED='\033[0;31m'
       GREEN='\033[0;32m'
       YELLOW='\033[1;33m'
       BLUE='\033[0;34m'
       NC='\033[0m' # No Color
   else
       RED=''
       GREEN=''
       YELLOW=''
       BLUE=''
       NC=''
   fi

   # -----------------------------------------------------------------------------
   # SEED DATA REFERENCE
   # -----------------------------------------------------------------------------
   # [Document available test data here]
   # Customers:
   #   - ID: xxx, Name: xxx, Status: xxx
   # Plans:
   #   - ID: xxx, Name: xxx, Rate: xxx
   # Quotas:
   #   - ID: xxx, Limit: xxx, Reset: xxx
   # -----------------------------------------------------------------------------

   # -----------------------------------------------------------------------------
   # TEST COUNTERS
   # -----------------------------------------------------------------------------
   TESTS_PASSED=0
   TESTS_FAILED=0
   TESTS_TOTAL=0

   # -----------------------------------------------------------------------------
   # HELPER FUNCTIONS
   # -----------------------------------------------------------------------------
   print_test_header() {
       echo ""
       echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
       echo -e "${BLUE}TEST: $1${NC}"
       echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
   }

   print_expected() {
       echo -e "${YELLOW}Expected: $1${NC}"
   }

   print_result() {
       echo -e "${GREEN}Response:${NC}"
   }

   # Check test result - called after each curl command
   # Usage: check_result "Test Name" "Expected Status" "$HTTP_CODE" "$BODY"
   check_result() {
       local test_name="$1"
       local expected_status="$2"
       local actual_status="$3"
       local body="$4"

       echo "$body"
       echo ""

       if [ "$actual_status" = "$expected_status" ]; then
           echo -e "${GREEN}✓ PASSED${NC} [HTTP Status: $actual_status]"
           TESTS_PASSED=$((TESTS_PASSED + 1))
       else
           echo -e "${RED}✗ FAILED${NC} [HTTP Status: $actual_status, Expected: $expected_status]"
           TESTS_FAILED=$((TESTS_FAILED + 1))
       fi
       echo ""
   }

   # -----------------------------------------------------------------------------
   # TEST CASES
   # -----------------------------------------------------------------------------
   ```

   **IMPORTANT**: Do NOT use `eval` or `run_test` wrapper function with complex quoting.
   Instead, use direct curl calls for each test case (see test case format below).

4. **Generate test cases for each acceptance criterion**

   For each identified acceptance criterion or API endpoint:

   a. **Happy path tests**:
   - Valid requests with expected successful responses
   - Test with available seed data

   b. **Validation error tests**:
   - Missing required fields
   - Invalid field formats
   - Out-of-range values

   c. **Edge case tests**:
   - Zero values (e.g., zero tokens, zero amount)
   - Empty strings or null values
   - Boundary conditions

   d. **Error scenario tests**:
   - Not found resources (invalid IDs)
   - Conflict scenarios
   - Unauthorized access (if applicable)

   **Test case format** (use direct curl calls, NOT eval-based wrapper):

   ```bash
   # -----------------------------------------------------------------------------
   # AC[N]: [Acceptance Criterion Description]
   # -----------------------------------------------------------------------------
   TEST_NAME="AC[N]: [Short Description]"
   EXPECTED="[Expected HTTP Status]"
   TESTS_TOTAL=$((TESTS_TOTAL + 1))
   print_test_header "$TEST_NAME"
   print_expected "HTTP $EXPECTED"
   print_result
   HTTP_CODE=$(curl -s -o /tmp/response.txt -w "%{http_code}" -X [METHOD] "${BASE_URL}/api/[endpoint]" \
       -H "Content-Type: application/json" \
       -m 10 \
       -d '{"[field]": "[value]"}')
   BODY=$(cat /tmp/response.txt)
   check_result "$TEST_NAME" "$EXPECTED" "$HTTP_CODE" "$BODY"
   ```

5. **Add edge case and negative tests**

   Include additional tests beyond acceptance criteria:

   ```bash
   # -----------------------------------------------------------------------------
   # EDGE CASES
   # -----------------------------------------------------------------------------

   # Edge Case: Zero tokens
   TEST_NAME="Edge Case: Zero Token Count"
   EXPECTED="[Expected Status]"
   TESTS_TOTAL=$((TESTS_TOTAL + 1))
   print_test_header "$TEST_NAME"
   print_expected "HTTP $EXPECTED"
   print_result
   HTTP_CODE=$(curl -s -o /tmp/response.txt -w "%{http_code}" -X POST "${BASE_URL}/api/[endpoint]" \
       -H "Content-Type: application/json" \
       -m 10 \
       -d '{"tokenCount": 0}')
   BODY=$(cat /tmp/response.txt)
   check_result "$TEST_NAME" "$EXPECTED" "$HTTP_CODE" "$BODY"

   # Edge Case: Missing required field
   TEST_NAME="Edge Case: Missing Required Field"
   EXPECTED="400"
   TESTS_TOTAL=$((TESTS_TOTAL + 1))
   print_test_header "$TEST_NAME"
   print_expected "HTTP $EXPECTED"
   print_result
   HTTP_CODE=$(curl -s -o /tmp/response.txt -w "%{http_code}" -X POST "${BASE_URL}/api/[endpoint]" \
       -H "Content-Type: application/json" \
       -m 10 \
       -d '{}')
   BODY=$(cat /tmp/response.txt)
   check_result "$TEST_NAME" "$EXPECTED" "$HTTP_CODE" "$BODY"

   # Edge Case: Invalid ID format
   TEST_NAME="Edge Case: Invalid ID Format"
   EXPECTED="404"
   TESTS_TOTAL=$((TESTS_TOTAL + 1))
   print_test_header "$TEST_NAME"
   print_expected "HTTP $EXPECTED"
   print_result
   HTTP_CODE=$(curl -s -o /tmp/response.txt -w "%{http_code}" -X GET "${BASE_URL}/api/[endpoint]/invalid-id" \
       -H "Content-Type: application/json" \
       -m 10)
   BODY=$(cat /tmp/response.txt)
   check_result "$TEST_NAME" "$EXPECTED" "$HTTP_CODE" "$BODY"
   ```

6. **Add cleanup and test summary footer**

   ```bash
   # -----------------------------------------------------------------------------
   # CLEANUP
   # -----------------------------------------------------------------------------
   rm -f /tmp/response.txt

   # -----------------------------------------------------------------------------
   # TEST SUMMARY
   # -----------------------------------------------------------------------------
   echo ""
   echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
   echo -e "${BLUE}TEST EXECUTION COMPLETE${NC}"
   echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
   echo ""
   echo "Base URL: ${BASE_URL}"
   echo "Finished at: $(date)"
   echo ""
   echo -e "Tests Passed: ${GREEN}${TESTS_PASSED}${NC}"
   echo -e "Tests Failed: ${RED}${TESTS_FAILED}${NC}"
   echo -e "Total Tests:  ${TESTS_TOTAL}"
   echo ""

   # Calculate pass rate
   if [ "$TESTS_TOTAL" -gt 0 ]; then
       PASS_RATE=$((TESTS_PASSED * 100 / TESTS_TOTAL))
       if [ "$TESTS_FAILED" -eq 0 ]; then
           echo -e "${GREEN}✓ All tests passed! (${PASS_RATE}%)${NC}"
       else
           echo -e "${RED}✗ Some tests failed (${PASS_RATE}% passed)${NC}"
       fi
   fi
   echo ""

   # Exit with error code if any tests failed
   if [ "$TESTS_FAILED" -gt 0 ]; then
       exit 1
   fi
   ```

7. **Write script to file and make executable**

   a. **Create scripts directory if needed**:
   - Ensure `scripts/` directory exists under the project root

   b. **Write the complete script**:
   - Save to `scripts/test-api.sh`

   c. **Make script executable**:
   - Run `chmod +x scripts/test-api.sh`

8. **Report generation summary**

   ```
   ✅ API test script generated and saved to `scripts/test-api.sh`

   📋 Test coverage:
   - Acceptance Criteria tests: [count]
   - Edge case tests: [count]
   - Negative tests: [count]
   - Total test cases: [count]

   📊 Seed data referenced:
   - Customers: [count]
   - Plans: [count]
   - [Other entities]: [count]

   📈 Test result tracking:
   - Pass/Fail counters with color-coded output
   - Pass rate percentage calculation
   - Exit code: 0 (all passed) or 1 (some failed)

   🚀 Usage:
      ./scripts/test-api.sh                    # Uses localhost:8080
      ./scripts/test-api.sh http://api:3000    # Custom base URL

   🔄 CI/CD Integration:
      The script returns exit code 1 if any tests fail,
      making it suitable for use in CI/CD pipelines.
   ```

**Output**

A self-contained, executable shell script (`scripts/test-api.sh`) with:
- cURL commands for all acceptance criteria scenarios
- Edge case and negative test scenarios
- Clear test output formatting
- No external dependencies (bash + curl only)
- Timeout protection (`-m 10`) on all requests
- HTTP status capture via `-o /tmp/response.txt -w "%{http_code}"`
- Documented seed data reference
- **Success/Failure counting** with pass rate calculation
- **Exit code** reflecting test results (0 = all passed, 1 = some failed)
- **Cleanup** of temp files after execution

**Script Requirements Checklist**

| Requirement | Implementation |
|-------------|----------------|
| Timeout | `-m 10` on every curl command |
| No external dependencies | Pure bash + curl, no jq/yq |
| HTTP status capture | `-o /tmp/response.txt -w "%{http_code}"` (NOT eval-based) |
| Test coverage | All ACs + edge cases + negative tests |
| Clear formatting | Header/expected/result for each test |
| Seed data reference | Commented at script top |
| Executable | `chmod +x` applied |
| Success/Failure counting | `TESTS_PASSED`, `TESTS_FAILED` counters with summary |
| Exit code | Exit 1 if any tests failed, exit 0 otherwise |
| Cleanup | `rm -f /tmp/response.txt` at end of script |

**Guardrails**

- Do NOT proceed without input context (API files or ACs)
- Do NOT use external tools like `jq`, `yq`, or `python` in the generated script
- Do NOT generate tests without proper timeout (`-m 10`)
- Do NOT skip edge cases (zero values, missing fields, invalid IDs)
- Do NOT hardcode test data without documenting it in the seed reference section
- **Do NOT use `eval` with `-w "\n%{http_code}"` — this causes shell quoting issues**
- **Do NOT use special characters (like `|`) in `-w` format strings**
- Always use `-o /tmp/response.txt -w "%{http_code}"` for HTTP status capture
- Always use `"${BASE_URL}"` variable for endpoint URLs
- Always make the script executable with `chmod +x`
- Always create `scripts/` directory if it does not exist
- Always cleanup temp files (`rm -f /tmp/response.txt`) at end of script
- Error messages in expected results MUST match actual API error messages
- Test data IDs MUST reference actual seed data or clearly indicate synthetic data

**Context Integrity Guardrails**:

- **MUST read ALL `@` referenced files completely** — do NOT skip or partially read
- **MUST search for seed data** in common locations before generating tests
- **Verify API endpoint paths** match actual implementation
- **Verify request body schemas** match actual DTOs/request objects
- **Verify error messages** match actual exception handler responses

**cURL Best Practices**

1. **Always include these flags**:
   - `-s` (silent mode, suppress progress)
   - `-m 10` (10 second timeout)
   - `-o /tmp/response.txt` (write body to temp file)
   - `-w "%{http_code}"` (capture only HTTP status code)
   - `-H "Content-Type: application/json"` (for JSON bodies)

2. **Request body formatting**:
   - Use `-d` with single-quoted JSON strings
   - Properly escape special characters in JSON
   - Use single quotes around JSON to avoid shell interpolation issues

3. **Variable interpolation**:
   - Use `"${BASE_URL}"` with quotes to handle special characters
   - Use `"${VARIABLE}"` syntax for all shell variables in URLs

**CRITICAL: Shell Quoting Pitfalls to Avoid**

When generating shell scripts with curl commands, avoid these common pitfalls:

1. **DO NOT use `eval` with `-w "\n%{http_code}"`**:
   - The `\n` escape sequence does not work reliably through `eval`
   - It often outputs literal `n` instead of a newline, breaking HTTP status extraction

2. **DO NOT use special characters in `-w` format string**:
   - Characters like `|` will be interpreted as shell pipe operators
   - Example: `-w '|||%{http_code}'` will fail with "syntax error near unexpected token `|'"

3. **CORRECT APPROACH - Use temp file for response body**:
   ```bash
   # Capture HTTP status code directly, write body to temp file
   HTTP_CODE=$(curl -s -o /tmp/response.txt -w "%{http_code}" -X POST "${BASE_URL}/api/endpoint" \
       -H "Content-Type: application/json" \
       -m 10 \
       -d '{"field": "value"}')
   BODY=$(cat /tmp/response.txt)
   ```

4. **Why this works**:
   - `-o /tmp/response.txt` writes the response body to a file (not stdout)
   - `-w "%{http_code}"` outputs ONLY the HTTP status code to stdout
   - No complex string parsing or newline handling needed
   - No `eval` required, avoiding all shell quoting issues

**Integration with SPDD Workflow**

This command is the **validation phase** of the SPDD workflow:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           SPDD Workflow                                  │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  Phase 1: /spdd-analysis                                                │
│  ┌────────────────────────────────────────────────────────────────┐    │
│  │ Business Requirement → Enriched Context                         │    │
│  └────────────────────────────────────────────────────────────────┘    │
│                              │                                          │
│                              ▼                                          │
│  Phase 2: /spdd-reasons-canvas                                         │
│  ┌────────────────────────────────────────────────────────────────┐    │
│  │ Enriched Context → REASONS Canvas Structured Prompt             │    │
│  └────────────────────────────────────────────────────────────────┘    │
│                              │                                          │
│                              ▼                                          │
│  Phase 3: /spdd-generate                                               │
│  ┌────────────────────────────────────────────────────────────────┐    │
│  │ Structured Prompt → Implementation Code                         │    │
│  └────────────────────────────────────────────────────────────────┘    │
│                              │                                          │
│                              ▼                                          │
│  Phase 4: /spdd-api-test    ← YOU ARE HERE                            │
│  ┌────────────────────────────────────────────────────────────────┐    │
│  │ Generated Code + ACs → API Test Script                          │    │
│  │                                                                 │    │
│  │ Output: scripts/test-api.sh                                     │    │
│  └────────────────────────────────────────────────────────────────┘    │
│                              │                                          │
│                              ▼                                          │
│  Phase 5: Validate & Iterate                                            │
│  ┌────────────────────────────────────────────────────────────────┐    │
│  │ Run tests → Identify issues → Update via /spdd-sync            │    │
│  └────────────────────────────────────────────────────────────────┘    │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

**Why This Phase Matters**

API testing validates that the generated implementation:
1. **Meets acceptance criteria**: Each AC is verified with actual HTTP requests
2. **Handles edge cases**: Zero values, missing fields, invalid data are tested
3. **Returns correct errors**: Error messages and status codes match specifications
4. **Works end-to-end**: Real HTTP calls verify the full request/response cycle

The self-contained nature ensures:
- Tests can run in any environment with just `bash` and `curl`
- No installation of additional tools required
- CI/CD pipelines can execute without setup complexity
- Results are immediately visible without parsing