import { useState } from 'react';
import { stores } from './stores';
import StoreSelector from './components/StoreSelector';
import StoreStatusDisplay from './components/StoreStatusDisplay';
import { useStoreStatus } from './hooks/useStoreStatus';

function App() {
    const [selectedStoreId, setSelectedStoreId] = useState(stores[0].id);
    const { storeStatus, loading, error, showStatus, fetchStoreStatus } = useStoreStatus();

    const handleStoreChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
        setSelectedStoreId(event.target.value);
    };

    const renderStoreStatus = () => {
        if (loading) {
            return <div className="p-4 text-blue-500">Loading...</div>;
        }
        if (error) {
            return <div className="p-4 text-red-500">Error: {error}</div>;
        }
        if (showStatus) {
            return <StoreStatusDisplay storeStatus={storeStatus} selectedStoreId={selectedStoreId} />;
        }
        return null;
    };

    return (
        <div className="flex flex-col items-center justify-center min-h-screen bg-gray-100">
            <header className="bg-white p-6 rounded-lg shadow-md text-center">
                <h1 className="text-2xl font-bold mb-4">Jironow</h1>
                <div className="mb-4">
                    <StoreSelector selectedStoreId={selectedStoreId} onStoreChange={handleStoreChange} />
                </div>
                <button
                    onClick={() => fetchStoreStatus(selectedStoreId)}
                    disabled={loading}
                    className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded disabled:opacity-50"
                >
                    {loading ? 'Loading...' : '営業状況を見る'}
                </button>
                {renderStoreStatus()}
            </header>
        </div>
    );
}

export default App;
