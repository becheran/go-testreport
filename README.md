# go-testreport

Generate a markdown test report from the go json test result
⏱️

# Test Report

Total 123 ✔️ Passed: 123 ⏩ Skipped: 2 ❌ Failed: 10 ⏱️ Duration: 123ms


<details>
    <summary>✔️ Package A - 10 sec  </summary>
    <blockquote>
        <details>
            <summary>✔️ Test/12lk             10s ⏱️</summary>
```
Foo
```
        </details>
    </blockquote>
</details>

<details>
    <summary>❌ Package B 10s ⏱️</summary>
    <blockquote>
        <details>
            <summary>⏩ Test/12lk             10s ⏱️</summary>
            <blockquote>
                foo
            </blockquote>
        </details>
    </blockquote>
</details>
