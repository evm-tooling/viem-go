#!/bin/bash
# =============================================================================
# Quick PR Script
# =============================================================================
# Creates a PR from staged files without switching your local branch.
# Automatically squash-merges and pulls the changes back to master.
#
# Usage: ./scripts/quick-pr.sh "pr-title"
#
# Requirements:
# - Must be on master branch locally
# - Must have staged files
# - gh CLI installed and authenticated
# =============================================================================

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
info() { echo -e "${BLUE}→${NC} $1"; }
success() { echo -e "${GREEN}✓${NC} $1"; }
warn() { echo -e "${YELLOW}⚠${NC} $1"; }
error() { echo -e "${RED}✗${NC} $1" >&2; exit 1; }

# =============================================================================
# Validation
# =============================================================================

# Check argument
if [ -z "$1" ]; then
    error "Usage: $0 \"pr-title\""
fi

PR_TITLE="$1"
# Convert title to branch-safe name (lowercase, replace spaces with dashes)
BRANCH_NAME=$(echo "$PR_TITLE" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9]/-/g' | sed 's/--*/-/g' | sed 's/^-//' | sed 's/-$//')

info "PR Title: $PR_TITLE"
info "Branch Name: $BRANCH_NAME"
echo ""

# Check we're on master
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [ "$CURRENT_BRANCH" != "master" ]; then
    error "You must be on the 'master' branch. Currently on: $CURRENT_BRANCH"
fi
success "On master branch"

# Check gh CLI is installed
if ! command -v gh &> /dev/null; then
    error "GitHub CLI (gh) is not installed. Install it with: brew install gh"
fi
success "GitHub CLI available"

# Check gh is authenticated
if ! gh auth status &> /dev/null; then
    error "GitHub CLI is not authenticated. Run: gh auth login"
fi
success "GitHub CLI authenticated"

# Check there are staged files
STAGED_FILES=$(git diff --cached --name-only)
if [ -z "$STAGED_FILES" ]; then
    error "No staged files. Stage files first with: git add <files>"
fi
success "Found staged files:"
echo "$STAGED_FILES" | sed 's/^/    /'
echo ""

# Check if branch already exists remotely
if git ls-remote --exit-code --heads origin "$BRANCH_NAME" &> /dev/null; then
    error "Branch '$BRANCH_NAME' already exists on remote. Choose a different PR title."
fi

# =============================================================================
# Create the PR without switching branches
# =============================================================================

info "Creating temporary worktree for branch operations..."

# Create a temporary directory for the worktree
WORKTREE_DIR=$(mktemp -d)
trap "rm -rf '$WORKTREE_DIR'" EXIT

# Create the new branch from current master HEAD
git branch "$BRANCH_NAME" HEAD
success "Created local branch: $BRANCH_NAME"

# Add worktree for the new branch
git worktree add "$WORKTREE_DIR" "$BRANCH_NAME" --quiet
success "Created temporary worktree"

# Copy staged changes to the worktree and commit
info "Applying staged changes to new branch..."

# Get the staged content as a patch and apply it in the worktree
git diff --cached --binary > "$WORKTREE_DIR/.staged.patch"

if [ -s "$WORKTREE_DIR/.staged.patch" ]; then
    (
        cd "$WORKTREE_DIR"
        git apply --index .staged.patch
        rm .staged.patch
        git commit -m "$PR_TITLE" --no-verify
    )
    success "Committed staged changes to $BRANCH_NAME"
else
    # Clean up on failure
    git worktree remove "$WORKTREE_DIR" --force 2>/dev/null || true
    git branch -D "$BRANCH_NAME" 2>/dev/null || true
    error "Failed to create patch from staged changes"
fi

# Push the branch
info "Pushing branch to origin..."
(cd "$WORKTREE_DIR" && git push -u origin "$BRANCH_NAME" --quiet --no-verify)
success "Pushed $BRANCH_NAME to origin"

# Clean up worktree (but keep the branch)
git worktree remove "$WORKTREE_DIR" --force
trap - EXIT  # Remove the trap since we cleaned up manually

# Create the PR
info "Creating pull request..."
PR_URL=$(gh pr create \
    --base master \
    --head "$BRANCH_NAME" \
    --title "$PR_TITLE" \
    --body "Auto-generated PR from staged changes." \
    --fill-first 2>&1 | tail -n1)
success "Created PR: $PR_URL"

# Wait a moment for GitHub to process
sleep 2

# Squash and merge
info "Squash-merging PR..."
gh pr merge "$BRANCH_NAME" --squash --delete-branch --admin
success "PR squash-merged and branch deleted"

# Pull changes to master
info "Pulling changes to master..."
git pull --quiet
success "master branch updated"

# Clean up local branch reference if it still exists
git branch -D "$BRANCH_NAME" 2>/dev/null || true

echo ""
echo -e "${GREEN}════════════════════════════════════════${NC}"
echo -e "${GREEN}✓ PR workflow complete!${NC}"
echo -e "${GREEN}════════════════════════════════════════${NC}"
echo ""
echo "Your staged changes are now merged into master."
echo "You stayed on master the entire time."
echo ""

# Note: staged files are still staged locally - user may want to unstage them
warn "Note: Your local staged files are still staged."
echo "      Run 'git reset' to unstage them, or they'll be included in your next commit."