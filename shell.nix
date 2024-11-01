{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    go
    postgresql
    gcc
    git
  ];

  shellHook = ''
    export GOPATH=$HOME/go
    export PATH=$GOPATH/bin:$PATH
    
    # Variables d'environnement pour PostgreSQL
    export PGDATA="$PWD/pgdata"
    export PGHOST="localhost"
    export PGPORT="5433"
    export PGUSER="soulsite"
    export PGPASSWORD="AtsiMemel6904**"
    export DATABASE_URL="postgresql://$PGUSER:$PGPASSWORD@$PGHOST:$PGPORT/soulsitedb"

    # Fonction pour vérifier si le port est disponible
    function is_port_available() {
      ! nc -z localhost $PGPORT
    }

    # Fonction pour vérifier si le serveur est en cours d'exécution
    function is_postgres_running() {
      pg_ctl status > /dev/null 2>&1
      return $?
    }

    # Fonction pour démarrer PostgreSQL
    function pg_start() {
      echo "Attempting to start PostgreSQL..."
      
      if ! is_port_available; then
        echo "Warning: Port $PGPORT is already in use. Stopping any existing PostgreSQL instance..."
        pg_ctl stop -D $PGDATA -m fast || true
        sleep 2
      fi

      if is_postgres_running; then
        echo "PostgreSQL is already running"
      else
        pg_ctl start -D $PGDATA -l postgres.log -o "-p $PGPORT"
        
        # Attendre que le serveur démarre
        for i in {1..30}; do
          if pg_isready -h localhost -p $PGPORT > /dev/null 2>&1; then
            echo "PostgreSQL is now running on port $PGPORT"
            return 0
          fi
          echo "Waiting for PostgreSQL to start... ($i/30)"
          sleep 1
        done
        
        echo "Failed to start PostgreSQL. Check postgres.log for details."
        return 1
      fi
    }

    # Fonction pour arrêter PostgreSQL
    function pg_stop() {
      if is_postgres_running; then
        echo "Stopping PostgreSQL..."
        pg_ctl stop -D $PGDATA -m fast
      else
        echo "PostgreSQL is not running"
      fi
    }

    # Fonction pour initialiser la base de données
    function init_database() {
      if [ ! -d $PGDATA ]; then
        echo "Initializing PostgreSQL database..."
        
        # Nettoyer les anciennes installations si nécessaire
        rm -rf $PGDATA postgres.log
        
        # Initialiser avec des paramètres explicites
        initdb --pgdata=$PGDATA --auth=trust --encoding=UTF8 --locale=C --username=postgres

        # Configuration de PostgreSQL
        echo "Configuring PostgreSQL..."
        cat >> $PGDATA/postgresql.conf <<EOF
listen_addresses = 'localhost'
port = $PGPORT
unix_socket_directories = '$PWD'
EOF

        cat >> $PGDATA/pg_hba.conf <<EOF
# IPv4 local connections:
host    all             all             127.0.0.1/32            trust
# IPv6 local connections:
host    all             all             ::1/128                 trust
# Unix domain socket
local   all             all                                     trust
EOF

        # Démarrer PostgreSQL
        echo "Starting PostgreSQL for initial setup..."
        if pg_start; then
          echo "Creating database role and database..."
          # Utiliser l'utilisateur postgres pour créer le rôle soulsite
          psql -p $PGPORT -U postgres -d postgres -c "CREATE ROLE soulsite WITH LOGIN SUPERUSER PASSWORD 'AtsiMemel6904**';"
          
          # Créer la base de données avec le nouvel utilisateur
          psql -p $PGPORT -U postgres -d postgres -c "CREATE DATABASE soulsitedb OWNER soulsite;"
          
          echo "Database initialization completed successfully!"
        else
          echo "Failed to initialize database. Check the logs."
          return 1
        fi
      else
        echo "PostgreSQL data directory already exists"
        pg_start
      fi
    }

    # Nettoyer l'environnement si nécessaire
    if [ -f $PGDATA/postmaster.pid ]; then
      echo "Cleaning up stale PID file..."
      rm -f $PGDATA/postmaster.pid
    fi

    # Initialiser et démarrer PostgreSQL
    init_database

    echo ""
    echo "PostgreSQL environment ready!"
    echo "Available commands:"
    echo "  pg_start  - Start PostgreSQL server"
    echo "  pg_stop   - Stop PostgreSQL server"
    echo "  psql \"postgresql://$PGUSER:$PGPASSWORD@$PGHOST:$PGPORT/soulsitedb\" - Connect to the database"
    echo ""

    # Afficher les informations de connexion
    echo "Database connection info:"
    echo "  Host: $PGHOST"
    echo "  Port: $PGPORT"
    echo "  Database: soulsitedb"
    echo "  User: $PGUSER"
    echo ""
  '';
}
