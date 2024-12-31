import { defineStore } from 'pinia';

interface Connection {
  id: number;
  name: string;
  host: string;
  username: string;
}

export const useConnectionStore = defineStore('connection', () => {
  const connections = ref<Connection[]>([]);
  const activeConnectionId = ref<number | null>(null);

  const addConnection = (connection: Connection) => {
    connections.value.push(connection);
    activeConnectionId.value = connection.id;
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

  return {
    connections,
    activeConnectionId,
    hasConnections,
    addConnection,
    removeConnection,
    setActiveConnection,
  };
});
