/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("7kff2zw80a7rmbu")

  // remove
  collection.schema.removeField("zhvihmzi")

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("7kff2zw80a7rmbu")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "zhvihmzi",
    "name": "settings",
    "type": "json",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {}
  }))

  return dao.saveCollection(collection)
})
