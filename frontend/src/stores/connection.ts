import { defineStore } from 'pinia';

interface Connection {
  id: number;
  connId: number;
  label: string;
  host: string;
  username: string;
  isConnecting?: boolean;
  connectionError?: boolean;
  errorMessage?: string;
  errorDetails?: string;
}

export const useConnectionStore = defineStore('connection', () => {
  const connections = ref<Connection[]>([]);
  const activeConnectionId = ref<number | null>(null);

  const getNextNumberForName = (baseName: string) => {
    const existingNumbers = connections.value
      .filter(conn => conn.label.startsWith(baseName))
      .map(conn => {
        const match = conn.label.match(new RegExp(`^${baseName}\\s*\\((\\d+)\\)$`));
        return match ? parseInt(match[1]) : 0;
      })
      .filter(num => !isNaN(num));

    if (existingNumbers.length === 0) return 1;
    return Math.max(...existingNumbers) + 1;
  };

  const addConnection = (connection: Connection) => {
    const baseName = connection.label;
    const existingSameName = connections.value.some(conn => conn.label === baseName);

    if (existingSameName) {
      const nextNumber = getNextNumberForName(baseName);
      connection.label = `${baseName} (${nextNumber})`;
    }

    connections.value.push({
      ...connection,
      isConnecting: true,
      connectionError: false,
      errorMessage: '',
    });
    activeConnectionId.value = connection.id;
  };

  const updateConnectionStatus = (id: number, status: Partial<Connection>) => {
    const connection = connections.value.find(c => c.id === id);
    if (connection) {
      Object.assign(connection, status);
    }
  };

  const removeConnection = (id: number) => {
    const index = connections.value.findIndex(c => c.id === id);
    if (index > -1) {
      connections.value.splice(index, 1);
      if (activeConnectionId.value === id) {
        activeConnectionId.value = connections.value[0]?.id ?? null;
      }
    }
  };

  const setActiveConnection = (id: number) => {
    activeConnectionId.value = id;
  };

  const hasConnections = computed(() => connections.value.length > 0);

  const activeConnection = computed(() => connections.value.find(c => c.id === activeConnectionId.value));

  return {
    connections,
    activeConnectionId,
    activeConnection,
    hasConnections,
    addConnection,
    removeConnection,
    setActiveConnection,
    updateConnectionStatus,
  };
});
