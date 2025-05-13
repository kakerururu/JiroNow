import './App.css';
import { useState } from 'react';
import { getJiroNowAPI } from './api/jiroNowAPI';
import type { GetStoreStatusResponseResponse } from './api/jiroNowAPI.schemas';
import { stores } from './stores';
import StoreSelector from './components/StoreSelector';
import StoreStatusDisplay from './components/StoreStatusDisplay';

function App() {
    const [selectedStoreId, setSelectedStoreId] = useState(stores[0].id);
    const [storeStatus, setStoreStatus] = useState<GetStoreStatusResponseResponse | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [showStatus, setShowStatus] = useState(false);

    const fetchStoreStatus = async () => {
        setLoading(true);
        setError(null);
        setShowStatus(false);
        try {
            const api = getJiroNowAPI();
            const response = await api.getStoreStatus({ xAccountId: selectedStoreId });
            setStoreStatus(response.data);
            setShowStatus(true);
        } catch (err) {
            setError('Failed to fetch store status');
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    const handleStoreChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
        setSelectedStoreId(event.target.value);
        setStoreStatus(null);
        setError(null);
        setShowStatus(false);
    };

    return (
        <div className="App">
            <header className="App-header">
                <h1>Jironow</h1>
                <StoreSelector selectedStoreId={selectedStoreId} onStoreChange={handleStoreChange} />
                <button onClick={fetchStoreStatus} disabled={loading}>
                    {loading ? 'Loading...' : '営業状況を見る'}
                </button>
                <StoreStatusDisplay
                    storeStatus={storeStatus}
                    loading={loading}
                    error={error}
                    showStatus={showStatus}
                    selectedStoreId={selectedStoreId}
                />
            </header>
        </div>
    );
}

export default App;
