import React from 'react';
import { stores } from '../stores';
import './StoreSelector.css';

interface StoreSelectorProps {
    selectedStoreId: string;
    onStoreChange: (event: React.ChangeEvent<HTMLSelectElement>) => void;
}

function StoreSelector({ selectedStoreId, onStoreChange }: StoreSelectorProps) {
    return (
        <div className="store-selector-container">
            <label htmlFor="store-select">店舗を選択して下さい</label>
            <select id="store-select" value={selectedStoreId} onChange={onStoreChange}>
                {stores.map(store => (
                    <option key={store.id} value={store.id}>
                        {store.name}
                    </option>
                ))}
            </select>
        </div>
    );
}

export default StoreSelector;
