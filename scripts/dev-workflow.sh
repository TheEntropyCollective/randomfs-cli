#!/bin/bash

# RandomFS Development Workflow Script
set -e

PROJECTS=("randomfs-core" "randomfs-cli" "randomfs-http" "randomfs-web")
CORE_PROJECT="randomfs-core"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to setup development environment with replace directives
setup_dev() {
    log "Setting up development environment with local replace directives..."
    
    for project in "${PROJECTS[@]}"; do
        if [ "$project" != "$CORE_PROJECT" ] && [ -d "../$project" ]; then
            cd "../$project"
            if [ -f "go.mod" ]; then
                # Add replace directive if not present
                if ! grep -q "replace github.com/TheEntropyCollective/randomfs-core" go.mod; then
                    echo "" >> go.mod
                    echo "replace github.com/TheEntropyCollective/randomfs-core => ../$CORE_PROJECT" >> go.mod
                    log "Added replace directive to $project"
                fi
                go mod tidy
                go mod vendor
            fi
        fi
    done
    log "Development environment setup complete"
}

# Function to build all projects
build_all() {
    log "Building all projects..."
    
    for project in "${PROJECTS[@]}"; do
        if [ -d "../$project" ]; then
            cd "../$project"
            if [ -f "go.mod" ]; then
                log "Building $project..."
                go build -o "$project" .
            fi
        fi
    done
    log "All projects built successfully"
}

# Function to test all projects
test_all() {
    log "Testing all projects..."
    
    for project in "${PROJECTS[@]}"; do
        if [ -d "../$project" ]; then
            cd "../$project"
            if [ -f "go.mod" ]; then
                log "Testing $project..."
                go test ./...
            fi
        fi
    done
    log "All tests passed"
}

# Function to commit and push changes
commit_and_push() {
    local message="$1"
    local project="$2"
    
    if [ -z "$project" ]; then
        error "Project name required for commit_and_push"
        return 1
    fi
    
    cd "../$project"
    git add .
    git commit -m "$message"
    git push origin main
    log "Committed and pushed changes to $project"
}

# Function to release a new version
release_version() {
    local version="$1"
    local project="$2"
    
    if [ -z "$version" ] || [ -z "$project" ]; then
        error "Version and project name required for release"
        return 1
    fi
    
    cd "../$project"
    git tag "v$version"
    git push origin "v$version"
    log "Released $project v$version"
}

# Function to update dependent projects to use published version
update_dependents() {
    local version="$1"
    
    if [ -z "$version" ]; then
        error "Version required for update_dependents"
        return 1
    fi
    
    for project in "${PROJECTS[@]}"; do
        if [ "$project" != "$CORE_PROJECT" ] && [ -d "../$project" ]; then
            cd "../$project"
            if [ -f "go.mod" ]; then
                log "Updating $project to use randomfs-core v$version..."
                
                # Remove replace directive
                sed -i '' '/replace github.com\/TheEntropyCollective\/randomfs-core/d' go.mod
                
                # Update require directive
                go mod edit -require="github.com/TheEntropyCollective/randomfs-core@v$version"
                
                go mod tidy
                go mod vendor
                
                commit_and_push "Update to randomfs-core v$version" "$project"
            fi
        fi
    done
}

# Function to show current status
status() {
    log "Current development status:"
    
    for project in "${PROJECTS[@]}"; do
        if [ -d "../$project" ]; then
            cd "../$project"
            if [ -f "go.mod" ]; then
                echo "  $project:"
                if grep -q "replace github.com/TheEntropyCollective/randomfs-core" go.mod; then
                    echo "    - Using local randomfs-core (dev mode)"
                else
                    echo "    - Using published randomfs-core"
                fi
                echo "    - Branch: $(git branch --show-current)"
                echo "    - Status: $(git status --porcelain | wc -l) changes"
            fi
        fi
    done
}

# Main script logic
case "${1:-help}" in
    "setup")
        setup_dev
        ;;
    "build")
        build_all
        ;;
    "test")
        test_all
        ;;
    "commit")
        commit_and_push "$2" "$3"
        ;;
    "release")
        release_version "$2" "$3"
        ;;
    "update-deps")
        update_dependents "$2"
        ;;
    "status")
        status
        ;;
    "full-release")
        if [ -z "$2" ]; then
            error "Version required for full-release"
            exit 1
        fi
        log "Performing full release process for version $2..."
        test_all
        commit_and_push "Release preparation for v$2" "$CORE_PROJECT"
        release_version "$2" "$CORE_PROJECT"
        update_dependents "$2"
        log "Full release process complete for v$2"
        ;;
    *)
        echo "RandomFS Development Workflow"
        echo ""
        echo "Usage: $0 <command> [args...]"
        echo ""
        echo "Commands:"
        echo "  setup                    Setup development environment with replace directives"
        echo "  build                    Build all projects"
        echo "  test                     Test all projects"
        echo "  commit <msg> <project>   Commit and push changes for a project"
        echo "  release <version> <project>  Release a new version of a project"
        echo "  update-deps <version>    Update all dependent projects to use published version"
        echo "  status                   Show current development status"
        echo "  full-release <version>   Complete release process (test, commit, release, update deps)"
        echo ""
        echo "Examples:"
        echo "  $0 setup"
        echo "  $0 build"
        echo "  $0 commit 'Add new feature' randomfs-core"
        echo "  $0 release 0.1.5 randomfs-core"
        echo "  $0 full-release 0.1.5"
        ;;
esac
