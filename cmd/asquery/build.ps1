$TARGETS_CONFS = ("windows","386"),("windows","amd64") 

Write-Host "[*] Starting asquery build" -ForegroundColor DarkGreen

# Time our execution
$sw = [System.Diagnostics.StopWatch]::startnew()

foreach ($i in $TARGETS_CONFS) {
    Write-Host "Building asquery for: " $i 
    $env:GOOS = $i[0]
    $env:GOARCH = $i[1]
    $current_dir = Split-Path (Get-Location) -Leaf
    $output_file_name = $current_dir + "_" +$i[0]+ "_" + $i[1]
    try {
        if ($i[0] -eq "windows") {
            $output_file_name = $output_file_name+".exe"
        }
        go build -o $output_file_name
        Write-Host "Finished building" -ForegroundColor Green
    }
    catch {
        Write-Host "Something failed" -ForegroundColor Red 
    }
}

Write-Host "[*] Build took $($sw.ElapsedMilliseconds) ms" -ForegroundColor DarkGreen



