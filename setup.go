import (
    "os"
    "os/exec"
    "fmt"
    "runtime"
    "strings"
    "path/filepath"
)

func runGlobalSetup() {
    execPath, err := os.Executable()
    if err != nil {
        fmt.Println("❌ Failed to locate the executable:", err)
        return
    }

    var destPath string
    switch runtime.GOOS {
    case "windows":
        destPath = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "Microsoft", "WindowsApps", "urlshort.exe")
    case "darwin", "linux":
        destPath = "/usr/local/bin/urlshort"
    default:
        fmt.Println("❌ Unsupported OS for global install.")
        return
    }

    if runtime.GOOS != "windows" && os.Geteuid() != 0 {
        fmt.Println("⚠️  Please run with sudo to install globally:")
        fmt.Println("   sudo ./urlshort --setup")
        return
    }

    input, err := os.ReadFile(execPath)
    if err != nil {
        fmt.Println("❌ Failed to read binary:", err)
        return
    }

    err = os.WriteFile(destPath, input, 0755)
    if err != nil {
        fmt.Println("❌ Failed to write binary to destination:", err)
        return
    }

    fmt.Println("✅ urlshort installed globally at:", destPath)
}
