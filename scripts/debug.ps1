[cmdletbinding()]
param(
  [Parameter()]
  [string]$Target = 'cmd/flfa/main.go',
  [int]$WaitSeconds = 30,
  [switch]$StopOnFailure,
  [switch]$PassThru,
  [Parameter(ValueFromRemainingArguments)]
  $Arguments = 'play'
)

begin {
  $RootFolder = Split-Path -Parent $PSScriptRoot
  $Target = Join-Path -Path $RootFolder -ChildPath $Target
  $DebugLogFilePath = "$PSScriptRoot/.debug.log"

  $DebugProcessParameters = @{
    FilePath               = 'dlv'
    ErrorAction            = 'Stop'
    PassThru               = $true
    WindowStyle            = 'Hidden'
    RedirectStandardOutput = $DebugLogFilePath
    ArgumentList           = "debug --headless --api-version=2 --listen=127.0.0.1:8181 $Target -- $Arguments"
  }

  if (Test-Path -Path $DebugLogFilePath) {
    Try {
      Remove-Item -Path $DebugLogFilePath -Force -ErrorAction Stop
    } Catch {
      $ActiveDelves = Get-Process -Name dlv*
      if ($ActiveDelves.Count -ne 0) {
        throw "Unable to delete ${DebugLogFilePath}: It might be because an active delve process is holding a lock on the file:`n$($ActiveDelves.Id)"
      } else {
        throw "Unable to delete ${DebugLogFilePath}: $_"
      }
    }
  }
}
process {
  $InformationPreference = 'Continue'
  Write-Information "Starting delve with:`n`t$($DebugProcessParameters.ArgumentList)"
  $Process = Start-Process @DebugProcessParameters
  Write-Information "Started delve (Proccess ID: $($Process.Id)), waiting up to $WaitSeconds seconds for it to report that it is listening"
  $TimeElapsed = 0
  while ($TimeElapsed -lt $WaitSeconds) {
    if ($TimeElapsed -eq 0) {
      $PercentComplete = 0
    } else {
      $PercentComplete = $TimeElapsed / $WaitSeconds
    }
    $TimeRemaining = $WaitSeconds - $TimeElapsed
    $StatusMessage = "Waiting up to $TimeRemaining more seconds for delve to report that it is listening"
    Write-Progress -Activity Waiting -Status $StatusMessage -PercentComplete $PercentComplete -CurrentOperation 'Waiting...'
    If (Select-String -Path $DebugLogFilePath -Pattern 'API server listening' -Quiet) {
      Write-Progress -Activity Waiting -Completed
      Write-Information (Get-Content -Path $DebugLogFilePath)
      If ($PassThru) {
        $Process
      }
      exit
    } elseif ($Process.HasExited) {
      Write-Progress -Activity Waiting -Completed
      Write-Information "Delve process exited unexpectedly with exit code $($Process.ExitCode)"
      exit
    } else {
      $TimeElapsed++
      Start-Sleep -Seconds 1
    }
  }
  If ($StopOnFailure) {
    Stop-Process $Process -Force
    Write-Error "Waited $WaitSeconds seconds but delve did not report that it was listening. Killed dlv process $($Process.Id)"
  } else {
    Write-Error "Waited $WaitSeconds seconds but delve did not report that it was listening. Process is still running as $($Process.Id)"
  }
}
end {}