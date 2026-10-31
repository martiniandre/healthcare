DELETE FROM medication_interactions WHERE catalog_version = 'sample-1';
DELETE FROM medication_catalog WHERE catalog_version = 'sample-1';
DELETE FROM medication_classes WHERE code IN ('penicillins', 'cephalosporins', 'nsaids', 'pyrazolones', 'coumarin_anticoagulants');
