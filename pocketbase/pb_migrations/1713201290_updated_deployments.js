/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("h2e1cdq94xgdclh")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "v8cb6gmi",
    "name": "project",
    "type": "relation",
    "required": true,
    "presentable": false,
    "unique": false,
    "options": {
      "collectionId": "7kff2zw80a7rmbu",
      "cascadeDelete": false,
      "minSelect": null,
      "maxSelect": 1,
      "displayFields": null
    }
  }))

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "drawo3ga",
    "name": "blueprint",
    "type": "relation",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "collectionId": "vs5gr49hpah1g9q",
      "cascadeDelete": false,
      "minSelect": null,
      "maxSelect": 1,
      "displayFields": null
    }
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("h2e1cdq94xgdclh")

  // remove
  collection.schema.removeField("v8cb6gmi")

  // remove
  collection.schema.removeField("drawo3ga")

  return dao.saveCollection(collection)
})
