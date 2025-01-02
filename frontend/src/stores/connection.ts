import { defineStore } from 'pinia';

interface Connection {
  id: number;
  name: string;
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

  const addConnection = (connection: Connection) => {
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
