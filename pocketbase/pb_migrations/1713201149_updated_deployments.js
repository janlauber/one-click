/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("h2e1cdq94xgdclh")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "yjdzrfgx",
    "name": "avatar",
    "type": "file",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "mimeTypes": [],
      "thumbs": [],
      "maxSelect": 1,
      "maxSize": 5242880,
      "protected": false
    }
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("h2e1cdq94xgdclh")

  // remove
  collection.schema.removeField("yjdzrfgx")

  return dao.saveCollection(collection)
})
