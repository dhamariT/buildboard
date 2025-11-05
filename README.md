# buildboard

A Rails 8.1 application for tracking builds and project metrics.

## Prerequisites

- Ruby 3.3+ (see `.ruby-version`)
- PostgreSQL 16
- Bundler

## Local Development Setup

### Option 1: Using Docker Compose (Recommended for Quick Start)

If you don't have PostgreSQL installed locally, use Docker Compose to run the entire stack:

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd buildboard
   ```

2. **Set up environment variables**
   ```bash
   cp example.env .env
   ```

   Update the `.env` file with your configuration. For Docker Compose, use the database URL:
   ```
   DATABASE_URL=postgres://postgres:buildboard123@db:5432/buildboard_development
   ```

3. **Get development credentials**

   Obtain `config/credentials/development.key` from your team or project administrator.

4. **Start the application**
   ```bash
   docker-compose up
   ```

   The application will be available at `http://localhost:3000`

### Option 2: Local Installation (Native)

If you prefer to run Rails directly on your machine:

1. **Install dependencies**
   ```bash
   bundle install
   ```

2. **Set up environment variables**
   ```bash
   cp example.env .env
   ```

   Update the `.env` file with your local database configuration:
   ```
   DATABASE_URL=postgres://postgres:yourpassword@localhost:5432/buildboard_development
   ```

   Make sure to set all required API keys and secrets:
   - `SECRET_KEY_BASE` (generate with `rails secret`)
   - `LOCKBOX_MASTER_KEY` (generate with `Lockbox.generate_key` in Rails console)
   - `MASTER_KEY` (generate with `rails secret`)
   - Add any required API keys (Loops, Grok, OpenAI, Slack)

3. **Get development credentials**

   Obtain `config/credentials/development.key` from your team or project administrator.

4. **Set up the database**
   ```bash
   bin/rails db:create
   bin/rails db:migrate
   bin/rails db:seed  # Optional: load seed data
   ```

5. **Start the development server**
   ```bash
   bin/dev
   ```

   This starts both the Rails server and the CSS watcher. The application will be available at `http://localhost:3000`

## Building for Production

To build the application for production deployment:

1. **Precompile assets**
   ```bash
   RAILS_ENV=production bin/rails assets:precompile
   ```

2. **Build Docker image**
   ```bash
   docker build -t buildboard:latest .
   ```

## Common Commands

- **Run tests**: `bin/rails test` or `bin/rails test:system`
- **Run linters**:
  - Ruby: `bundle exec rubocop`
  - ERB: `bundle exec erblint --lint-all`
- **Database console**: `bin/rails dbconsole`
- **Rails console**: `bin/rails console`

## Tech Stack

- Rails 8.1.1
- PostgreSQL 16
- Hotwire (Turbo + Stimulus)
- Dart Sass for CSS
- Solid Queue, Solid Cache, Solid Cable for background jobs and caching
